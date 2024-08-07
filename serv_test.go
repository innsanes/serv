package serv_test

import (
	"github.com/stretchr/testify/assert"
	"serv"
	"testing"
	"time"
)

type ServiceSuccess struct {
	t           *testing.T
	serve       bool
	beforeServe bool
	afterServe  bool
	beforeStop  bool
	serv.Service
}

func (s *ServiceSuccess) BeforeServe() error {
	s.beforeServe = true
	assert.True(s.t, s.beforeServe)
	assert.False(s.t, s.serve)
	assert.False(s.t, s.afterServe)
	assert.False(s.t, s.beforeStop)
	return nil
}

func (s *ServiceSuccess) Serve() error {
	s.serve = true
	assert.True(s.t, s.beforeServe)
	assert.True(s.t, s.serve)
	assert.False(s.t, s.afterServe)
	assert.False(s.t, s.beforeStop)
	return nil
}

func (s *ServiceSuccess) AfterServe() error {
	s.afterServe = true
	assert.True(s.t, s.beforeServe)
	assert.True(s.t, s.serve)
	assert.True(s.t, s.afterServe)
	assert.False(s.t, s.beforeStop)
	return nil
}

func (s *ServiceSuccess) BeforeStop() error {
	s.beforeStop = true
	assert.True(s.t, s.beforeServe)
	assert.True(s.t, s.serve)
	assert.True(s.t, s.afterServe)
	assert.True(s.t, s.beforeStop)
	return nil
}

func TestServ(t *testing.T) {
	server := serv.New()
	go func() {
		time.Sleep(1 * time.Second)
		server.ForceStop()
	}()

	beforeServe := false
	afterServe := false
	beforeStop := false
	server.RegisterBeforeServe(func() error {
		beforeServe = true
		assert.True(t, beforeServe)
		assert.False(t, afterServe)
		assert.False(t, beforeStop)
		return nil
	})
	server.RegisterAfterServe(func() error {
		afterServe = true
		assert.True(t, beforeServe)
		assert.True(t, afterServe)
		assert.False(t, beforeStop)
		return nil
	})
	server.RegisterBeforeStop(func() error {
		beforeStop = true
		assert.True(t, beforeServe)
		assert.True(t, afterServe)
		assert.True(t, beforeStop)
		return nil
	})
	service := &ServiceSuccess{t: t}
	server.Serve(service, &ServiceSuccess{t: t})
	assert.True(t, service.beforeServe)
	assert.True(t, service.serve)
	assert.True(t, service.afterServe)
	assert.True(t, service.beforeStop)
}

type ServicePanic struct {
	t     *testing.T
	value bool
	serv.Service
}

func (s *ServicePanic) BeforeServe() error {
	s.value = true
	panic("BeforeServe")
}

func TestServPanic(t *testing.T) {
	server := serv.New()
	go func() {
		// 在启动中产生 panic 会导致服务终止退出 所以不会指定这段逻辑
		// 如果出现 panic 则说明服务没有终止
		time.Sleep(1 * time.Second)
		panic(1)
	}()

	service := &ServicePanic{t: t}
	server.Serve(service)
	assert.True(t, service.value)
}
