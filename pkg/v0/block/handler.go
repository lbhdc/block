package block

import (
	"context"
	"fmt"
	pb "github.com/lbhdc/block/api/v0/net/http"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"sync"
)

type HandlerPlugin struct {
	Plugin
	Addr   string
	Port   uint16
	conn   *grpc.ClientConn
	client pb.HandlerClient
	ctx    context.Context
}

func NewHandlerPlugin(ctx context.Context, cfg HandlerConfig) *HandlerPlugin {
	hp := &HandlerPlugin{
		Plugin: NewPlugin(cfg.Entrypoint),
		Addr:   cfg.Addr,
		Port:   cfg.Port,
		ctx:    ctx,
	}
	return hp
}

func (hp *HandlerPlugin) Connect() (err error) {
	hp.conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", hp.Port), grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Error("grpc.Dial")
		return err
	}
	hp.client = pb.NewHandlerClient(hp.conn)
	return nil
}

var onlyOnce sync.Once

func (hp *HandlerPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	onlyOnce.Do(func() {
		if err := hp.Connect(); err != nil {
			log.WithError(err).Fatalln("hp.Connect")
		}
	})
	var body []byte
	if r.Body != nil {
		var err error
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("ioutil.ReadAll")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
	var header []*pb.Header
	if len(r.Header) != 0 {
		for key, vals := range r.Header {
			for _, val := range vals {
				header = append(header, &pb.Header{Key: key, Value: val})
			}
		}
	}

	res, err := hp.client.Handle(hp.ctx, &pb.Request{
		Path:   r.URL.String(),
		Method: r.Method,
		Header: header,
		Body:   body,
	})
	if err != nil {
		log.WithError(err).Error("hp.client.Handle")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(int(res.Code))
	if _, err = w.Write(res.Body); err != nil {
		log.WithError(err).Error("w.Write")
	}
}
