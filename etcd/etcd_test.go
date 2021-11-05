package etcd_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ETCDTestSuite struct {
	suite.Suite
	cli *clientv3.Client
}

func TestETCDTestSuite(t *testing.T) {
	suite.Run(t, new(ETCDTestSuite))
}

func (s *ETCDTestSuite) SetupSuite() {
	var err error
	s.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:12379", "http://localhost:22379", "http://localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	s.NoError(err)
}

func (s *ETCDTestSuite) TearDownSuite() {
	s.cli.Close()
}

func (s *ETCDTestSuite) TestPutGetDeleteWithoutOption() {
	key, val := "/test/etcd/abc3", "ABC3"
	putRes, err := s.cli.Put(context.Background(), key, val)
	s.NoError(err)
	s.NotZero(putRes)
	s.Zero(putRes.PrevKv)
	getRes, err := s.cli.Get(context.Background(), key)
	s.NoError(err)
	for _, v := range getRes.Kvs {
		s.Equal(key, string(v.Key))
		s.Equal(val, string(v.Value))
	}
	_, err = s.cli.Delete(context.Background(), key)
	s.NoError(err)
}

func (s *ETCDTestSuite) TestGetDeleteWithPrefixOption() {
	prefix := "/test/etcd/"
	keys := []string{prefix + "abc1", prefix + "abc2", prefix + "abc3"}
	vals := []string{"ABC1", "ABC2", "ABC3"}
	for i, key := range keys {
		putRes, err := s.cli.Put(context.Background(), key, vals[i])
		s.NoError(err)
		s.NotZero(putRes)
	}

	getRes, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	s.NoError(err)

	for i, v := range getRes.Kvs {
		s.Equal(keys[i], string(v.Key))
		s.Equal(vals[i], string(v.Value))
	}
	_, err = s.cli.Delete(context.Background(), prefix, clientv3.WithPrefix())
	s.NoError(err)
}

func (s *ETCDTestSuite) TestKV() {
	kv := clientv3.NewKV(s.cli)
	key, val := "/cron/jobs/job1", "hello"

	_, err := kv.Put(context.Background(), key, val)
	s.NoError(err)

	s.Run("WithPrevKV", func() {
		resp, err := kv.Put(context.Background(), key, "world", clientv3.WithPrevKV())
		s.NoError(err)
		s.Equal(val, string(resp.PrevKv.Value))
		_, err = kv.Delete(context.Background(), key)
		s.NoError(err)
	})
}

func (s *ETCDTestSuite) TestWithLeaseOption() {
	lease := clientv3.NewLease(s.cli)
	ttl := 2
	resp, err := lease.Grant(context.Background(), int64(ttl))
	s.NoError(err)
	s.Equal(int64(ttl), resp.TTL)

	key, val := "/test/lease/abc", "ABC3"
	putRes, err := s.cli.Put(context.Background(), key, val, clientv3.WithLease(resp.ID))
	s.NoError(err)
	s.NotZero(putRes)

	time.Sleep(3 * time.Second)

	getRes, err := s.cli.Get(context.Background(), key)
	s.NoError(err)
	s.NotZero(getRes)
	s.Zero(getRes.Kvs)
	s.Equal(int64(0), getRes.Count)

	delResp, err := s.cli.Delete(context.Background(), key)
	s.NoError(err)
	s.NotZero(delResp)
	s.Zero(delResp.PrevKvs)
}

func (s *ETCDTestSuite) TestWithLeaseOptionAndKeepAlive() {

	// create Lease
	lease := clientv3.NewLease(s.cli)
	ttl := 2
	// Grant ttl to Lease
	resp, err := lease.Grant(context.Background(), int64(ttl))
	s.NoError(err)
	s.Equal(int64(ttl), resp.TTL)

	key, val := "/test/lease/abc", "ABC3"
	putRes, err := s.cli.Put(context.Background(), key, val, clientv3.WithLease(resp.ID))
	s.NoError(err)
	s.NotZero(putRes)

	// auto extend ttl
	ctx, cancleFunc := context.WithCancel(context.Background())
	keepChan, err := lease.KeepAlive(ctx, resp.ID)
	s.NoError(err)

	go func() {
		//  get the result of extending ttl
		keepResp := <-keepChan
		s.NotZero(keepResp)
		s.Equal(resp.ID, keepResp.ID)

		//  get the result of extending ttl
		keepResp = <-keepChan
		s.NotZero(keepResp)
		s.Equal(resp.ID, keepResp.ID)

		// stop extending ttl
		cancleFunc()

		// can't get result any more
		keepResp = <-keepChan
		s.Zero(keepResp)
	}()

	// now ttl is 2 + 2
	time.Sleep(4 * time.Second)

	getRes, err := s.cli.Get(context.Background(), key)
	s.NoError(err)
	s.NotZero(getRes)
	s.Zero(getRes.Kvs)
	s.Equal(int64(0), getRes.Count)

	delResp, err := s.cli.Delete(context.Background(), key)
	s.NoError(err)
	s.NotZero(delResp)
	s.Zero(delResp.PrevKvs)
}

func (s *ETCDTestSuite) TestWatch() {
	prefix := "/test/watch/"
	keys := []string{prefix + "key", prefix + "key1"}
	vals := []string{"val", "val1"}
	watcher := clientv3.NewWatcher(s.cli)
	rch := watcher.Watch(context.Background(), prefix, clientv3.WithPrefix())

	go func() {
		for i, key := range keys {
			_, err := s.cli.Put(context.Background(), key, vals[i])
			s.NoError(err)
			time.Sleep(500 * time.Millisecond)
		}
		watcher.Close()
	}()

	i := 0
	for wresp := range rch {
		for _, ev := range wresp.Events {
			s.Equal("PUT", ev.Type.String())
			s.Equal(keys[i], string(ev.Kv.Key))
			s.Equal(vals[i], string(ev.Kv.Value))
			i++
		}
	}

	_, err := s.cli.Delete(context.Background(), prefix, clientv3.WithPrefix())
	s.NoError(err)

}
