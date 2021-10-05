package receivers

import (
	"S3Replicator/queue"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7/pkg/notification"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

const HTTP_TYPE_NAME = "http"

type HTTPCommandLineConfig struct {
	Host string `long:"host" description:"Listen host" required:"true" env:"LISTEN" default:"127.0.0.1"`
	Port uint16 `long:"port" description:"Listen Port" required:"true" env:"PORT" default:"8080"`
}

type HTTPEventHandler struct {
	Config              HTTPCommandLineConfig
	NotificationChannel *queue.EventChannelQueue
	srv                 *http.Server
}

func (H *HTTPEventHandler) Close() error {
	return H.srv.Close()
}

func (H *HTTPEventHandler) GetHandlerName() string {
	return HTTP_TYPE_NAME
}

func (H *HTTPEventHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	log.Debug("Received webhook request !", r.Method, r.RequestURI)
	bytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	ev := notification.Info{}
	err := json.Unmarshal(bytes, &ev)
	if err != nil {
		log.Error("ERR: ", err)
	}

	//For each record in the message produce the event
	for _, record := range ev.Records {
		H.NotificationChannel.Queue(record)
		log.Debug("Record queued", record)
		log.Debug("queue size: ", H.NotificationChannel.Size())
	}

	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (H *HTTPEventHandler) Init() error {
	log.Info("Initializing http server...")
	//Do init !
	router := mux.NewRouter()

	router.Methods("POST", "GET").Path("/webhook/event").HandlerFunc(H.handleWebhook)

	listenAdd := H.Config.Host + ":" + strconv.Itoa(int(H.Config.Port))
	H.srv = &http.Server{
		Handler: router,
		Addr:    listenAdd,
	}

	router.NotFoundHandler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Warn("Not found...", request.Method, request.RequestURI)
	})

	log.Info("Waiting for events on " + listenAdd + "/webhook/event ! ")

	go H.runSrv()
	return nil
}

func (H *HTTPEventHandler) runSrv() {
	log.Error(H.srv.ListenAndServe())
}
