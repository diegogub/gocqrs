package gocqrs

import (
	"encoding/json"
	"errors"
	"github.com/diegogub/lib"
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"strconv"
	"sync"
)

var (
	InvalidEntityError = errors.New("Invalid entity")
)

const (
	EventTypeHeader     = "X-Event"
	EntityVersionHeader = "X-LockVersion"
	CreateEventHeader   = "X-Create"
	EntityHeader        = "X-Entity"
	EventIDHeader       = "X-EventID"
)

var runningApp *App

type App struct {
	lock    sync.Mutex
	Version string `json:"version"`
	Name    string `json:"name"`
	Port    string `json:"port"`

	Entities map[string]*EntityConf `json:"entities"`
	Store    EventStore             `json:"-"`
	Router   *gin.Engine
}

func (app *App) String() string {
	b, _ := json.Marshal(app)
	return string(b)
}

func NewApp(store EventStore) *App {
	var app App
	app.Entities = make(map[string]*EntityConf)
	app.Router = gin.New()
	app.Store = store
	return &app
}

func (app *App) RegisterEntity(e *EntityConf) *App {
	_, has := app.Entities[e.Name]
	if !has {
		app.Entities[e.Name] = e
	} else {
		log.Fatal("Entity already added")
	}
	return app
}

func (app *App) HandleEvent(entityName, id string, ev Eventer, versionLock uint64) (string, uint64, error) {
	var err error
	app.lock.Lock()
	defer app.lock.Unlock()

	econf, ok := app.Entities[entityName]
	if !ok {
		return "", 0, InvalidEntityError
	}

	// look for entity events, TODO eventstore should cache streams
	stream := entityName + "-" + id
	ch := app.Store.Range(stream)
	entity, err := econf.Aggregate(id, ch)
	if err != nil {
		return "", 0, err
	}

	h, has := econf.EventHandlers[ev.GetType()]
	if !has {
		return "", 0, errors.New("Invalid handler for event:" + ev.GetType())
	}

	// handler event
	opt, err := h.Handle(ev, entity)
	if err != nil {
		return "", 0, err
	}

	//validate entity
	for n, v := range econf.Validators {
		err = v.Validate(*entity)
		if err != nil {
			return "", 0, errors.New("Failed validation: " + n + " - " + err.Error())
		}
	}

	version, err := app.Store.Store(ev, opt)
	return entity.ID, version, err
}

// Start app
func (app *App) Run(port string) error {
	app.Router.POST("/event/:entity", HTTPEventHandler)
	app.Router.GET("/docs", DocHandler)
	runningApp = app
	return runningApp.Router.Run(port)
}

func HTTPEventHandler(c *gin.Context) {
	e := c.Param("entity")
	data := make(map[string]interface{})

	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": err.Error()})
		return
	}

	eVersion := c.Request.Header.Get(EntityVersionHeader)
	v, _ := strconv.ParseUint(eVersion, 10, 64)

	// get event type
	eType := c.Request.Header.Get(EventTypeHeader)

	// get entity id
	enID := c.Request.Header.Get(EntityHeader)
	if enID == "" {
		enID = lib.NewShortId("")
	}

	// get event id
	eID := c.Request.Header.Get(EventIDHeader)

	event := NewEvent(eID, eType, data)
	event.Entity = e
	event.EntityID = enID

	// create event
	id, version, err := runningApp.HandleEvent(event.Entity, event.EntityID, event, v)
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(201, map[string]interface{}{"entity": e, "entity-id": enID, "version": version, "event-id": id})
	return
}

func DocHandler(c *gin.Context) {
	c.JSON(200, GenerateDocs(runningApp))
}
