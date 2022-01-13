package threads

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type TestRunnableThread struct {
	TestID      string
	TStartError error
	started     bool
}

func (trt *TestRunnableThread) GetID() string                      { return trt.TestID }
func (trt *TestRunnableThread) SetChannel(chan event.GenericEvent) {}
func (trt *TestRunnableThread) Start(context.Context, logr.Logger) error {
	if trt.TStartError == nil {
		trt.started = true
	}
	return trt.TStartError
}
func (trt *TestRunnableThread) Stop()           { trt.started = false }
func (trt *TestRunnableThread) IsStarted() bool { return trt.started }

func TestManager_RunThread(t *testing.T) {
	type fields struct {
		channel chan event.GenericEvent
		threads map[string]RunnableThread
	}
	type args struct {
		ctx    context.Context
		key    string
		thread RunnableThread
		log    logr.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Runs a RunnableThread",
			fields: fields{
				channel: make(chan event.GenericEvent),
				threads: map[string]RunnableThread{},
			},
			args: args{
				ctx: context.TODO(),
				thread: &TestRunnableThread{
					TestID:      "test",
					TStartError: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "Returns an error",
			fields: fields{
				channel: make(chan event.GenericEvent),
				threads: map[string]RunnableThread{},
			},
			args: args{
				ctx: context.TODO(),
				key: "key",
				thread: &TestRunnableThread{
					TestID:      "test",
					TStartError: errors.New("error"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := &Manager{
				channel: tt.fields.channel,
				threads: tt.fields.threads,
			}

			if err := mgr.RunThread(tt.args.ctx, tt.args.key, tt.args.thread, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("Manager.RunThread() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if _, ok := mgr.threads[tt.args.key]; !ok {
					t.Errorf("Manager.RunThread() RunnableThread not found in manager map")
				}
			}
		})
	}
}

func TestManager_StopThread(t *testing.T) {
	type fields struct {
		channel chan event.GenericEvent
		threads map[string]RunnableThread
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Stops a RunnableThread",
			fields: fields{
				channel: make(chan event.GenericEvent),
				threads: map[string]RunnableThread{
					"key": &TestRunnableThread{
						TestID:  "test",
						started: true,
					}},
			},
			args: args{key: "key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := &Manager{
				channel: tt.fields.channel,
				threads: tt.fields.threads,
			}
			mgr.StopThread(tt.args.key)
			if _, ok := mgr.threads[tt.args.key]; ok {
				t.Errorf("Manager.RunThread() RunnableThread should have been deleted from manager")
			}
		})
	}
}

func TestManager_GetChannel(t *testing.T) {

	t.Run("Returns the channel", func(t *testing.T) {
		ch := make(chan event.GenericEvent)
		mgr := &Manager{
			channel: ch,
			threads: map[string]RunnableThread{},
		}
		msg := event.GenericEvent{Object: &corev1.Namespace{}}
		go func() {
			ch <- msg
		}()

		if got := <-mgr.GetChannel(); !reflect.DeepEqual(got, msg) {
			t.Errorf("Manager.GetChannel() = %v, want %v", got, ch)
		}
	})
}

func TestManager_ReconcileThreads(t *testing.T) {
	type fields struct {
		channel chan event.GenericEvent
		threads map[string]RunnableThread
	}
	type args struct {
		ctx      context.Context
		instance client.Object
		threads  []RunnableThread
		log      logr.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]RunnableThread
		wantErr bool
	}{
		{
			name: "Adds new threads to the manager and runs them",
			fields: fields{
				channel: make(chan event.GenericEvent),
				threads: map[string]RunnableThread{},
			},
			args: args{
				instance: &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc", Namespace: "ns"}},
				threads: []RunnableThread{
					&TestRunnableThread{TestID: "t1", started: false},
					&TestRunnableThread{TestID: "t2", started: false},
					&TestRunnableThread{TestID: "t3", started: false},
				},
			},
			want: map[string]RunnableThread{
				"ns/svc_t1": &TestRunnableThread{TestID: "t1", started: true},
				"ns/svc_t2": &TestRunnableThread{TestID: "t2", started: true},
				"ns/svc_t3": &TestRunnableThread{TestID: "t3", started: true},
			},
			wantErr: false,
		},
		{
			name: "Returns error if one of the threads fail",
			fields: fields{
				channel: make(chan event.GenericEvent),
				threads: map[string]RunnableThread{},
			},
			args: args{
				instance: &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc", Namespace: "ns"}},
				threads: []RunnableThread{
					&TestRunnableThread{TestID: "t1", started: false, TStartError: errors.New("error")},
				},
			},
			want: map[string]RunnableThread{
				"ns/svc_t1": &TestRunnableThread{TestID: "t1", started: false, TStartError: errors.New("error")},
			},
			wantErr: true,
		},
		{
			name: "Starts a thread that previously failed",
			fields: fields{
				channel: make(chan event.GenericEvent),
				threads: map[string]RunnableThread{
					"ns/svc_t1": &TestRunnableThread{TestID: "t1", started: false},
				},
			},
			args: args{
				instance: &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc", Namespace: "ns"}},
				threads: []RunnableThread{
					&TestRunnableThread{TestID: "t1", started: false},
				},
			},
			want: map[string]RunnableThread{
				"ns/svc_t1": &TestRunnableThread{TestID: "t1", started: true},
			},
			wantErr: false,
		},
		{
			name: "Stops threads that should not run anymore",
			fields: fields{
				channel: make(chan event.GenericEvent),
				threads: map[string]RunnableThread{
					"ns/svc_t1": &TestRunnableThread{TestID: "t1", started: true},
					"ns/svc_t2": &TestRunnableThread{TestID: "t2", started: true},
					"ns/svc_t3": &TestRunnableThread{TestID: "t3", started: true},
				},
			},
			args: args{
				instance: &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc", Namespace: "ns"}},
				threads: []RunnableThread{
					&TestRunnableThread{TestID: "t1"},
				},
			},
			want: map[string]RunnableThread{
				"ns/svc_t1": &TestRunnableThread{TestID: "t1", started: true},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := &Manager{
				channel: tt.fields.channel,
				threads: tt.fields.threads,
			}
			if err := mgr.ReconcileThreads(tt.args.ctx, tt.args.instance, tt.args.threads, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("Manager.ReconcileThreads() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := deep.Equal(mgr.threads, tt.want); len(diff) > 0 {
				t.Errorf("Manager.ReconcileThreads() = diff %v", diff)
			}
		})
	}
}

func TestManager_CleanupThreads(t *testing.T) {
	type fields struct {
		channel chan event.GenericEvent
		threads map[string]RunnableThread
	}
	type args struct {
		instance client.Object
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]RunnableThread
	}{
		{
			name: "Returns a func that stops and cleans all threads belonging to an instance when invoked",
			fields: fields{
				threads: map[string]RunnableThread{
					"ns/svc_t1":   &TestRunnableThread{TestID: "t1", started: true},
					"ns/svc_t2":   &TestRunnableThread{TestID: "t2", started: true},
					"ns/other_t1": &TestRunnableThread{TestID: "t1", started: true},
				},
			},
			args: args{
				instance: &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc", Namespace: "ns"}},
			},
			want: map[string]RunnableThread{
				"ns/other_t1": &TestRunnableThread{TestID: "t1", started: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := &Manager{
				channel: tt.fields.channel,
				threads: tt.fields.threads,
			}
			mgr.CleanupThreads(tt.args.instance)()
			if diff := deep.Equal(mgr.threads, tt.want); len(diff) > 0 {
				t.Errorf("Manager.CleanupThreads() = diff %v", diff)
			}
		})
	}
}

func Test_prefix(t *testing.T) {
	type args struct {
		o client.Object
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Resturns the prefix for threads belonging to this instance",
			args: args{
				o: &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc", Namespace: "ns"}},
			},
			want: "ns/svc_",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prefix(tt.args.o); got != tt.want {
				t.Errorf("prefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
