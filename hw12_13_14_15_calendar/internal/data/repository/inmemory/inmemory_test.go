package inmemory

import (
	"context"
	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/internal/domain/entities"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var (
	testid    = "00112233-4455-6677-8899-aabbccddeeff"
	testdt, _ = time.Parse(entities.LayoutISO, "2020-06-25 14:14:14")
	testEvent = entities.Event{
		ID:         testid,
		Title:      "title",
		DateTime:   testdt,
		Duration:   360000000000,
		Text:       "text",
		UserID:     testid,
		TimeNotify: 360000000000,
	}
)

func TestInMemoryEventRepo(t *testing.T) {
	repo, err := NewInMemoryEventRepo(NewMapRepo(), nil)
	require.Nil(t, err)
	ctx := context.Background()
	t.Run("add good", func(t *testing.T) {
		id, err := repo.Add(ctx, testEvent)
		require.NotEqual(t, "", id, "generated id must be not empty")
		require.Nil(t, err)
		require.Equal(t, len(repo.m.events), 1)
		require.Equal(t, len(repo.m.users), 1)
		require.Equal(t, len(repo.m.users[testEvent.UserID]), 1)

	})
	t.Run("add bad", func(t *testing.T) {
		repo.m.Clear()
		outsideAddMapRepo(repo.m, testEvent)

		id, err := repo.Add(ctx, testEvent)
		require.Equal(t, "", id, "generated id must be empty")
		require.Truef(t, errors.Is(err, entities.ErrDateBusy), "return error must be: %q", entities.ErrDateBusy)
	})
	t.Run("get good", func(t *testing.T) {
		repo.m.Clear()
		outsideAddMapRepo(repo.m, testEvent)

		e, err := repo.GetByID(ctx, testid, testid)

		require.Nil(t, err)
		require.Equal(t, &testEvent, e, "event not expected")
	})
	t.Run("get user not found", func(t *testing.T) {
		repo.m.Clear()
		outsideAddMapRepo(repo.m, testEvent)

		e, err := repo.GetByID(ctx, "not id", testid)

		require.Nil(t, e)
		require.Truef(t, errors.Is(err, entities.ErrUnknownUser), "return error must be: %q", entities.ErrEventNotFound)
	})
	t.Run("get event not found", func(t *testing.T) {
		repo.m.Clear()
		outsideAddMapRepo(repo.m, testEvent)
		e, err := repo.GetByID(ctx, testid, "not id")

		require.Nil(t, e)
		require.Truef(t, errors.Is(err, entities.ErrEventNotFound), "return error must be: %q", entities.ErrEventNotFound)
	})
	t.Run("get by date good", func(t *testing.T) {
		repo.m.Clear()
		outsideAddMapRepo(repo.m, testEvent)

		e, err := repo.GetByDate(ctx, testid, testdt)

		require.Nil(t, err)
		var es []*entities.Event
		es = append(es, &testEvent)
		require.Equal(t, es, e, "events array not expected")
	})
	t.Run("get by date bad", func(t *testing.T) {
		repo.m.Clear()
		outsideAddMapRepo(repo.m, testEvent)
		e, err := repo.GetByDate(ctx, testid, time.Now())

		require.Nil(t, e)
		require.Truef(t, errors.Is(err, entities.ErrEventNotFound), "return error must be: %q", entities.ErrEventNotFound)
	})
	t.Run("get by period good", func(t *testing.T) {
		repo.m.Clear()

		testEvent2 := entities.Event{
			ID:         "11112233-4455-6677-8899-aabbccddeeff",
			Title:      "title",
			DateTime:   testdt.AddDate(0, 0, 5),
			Duration:   360000000000,
			Text:       "text",
			UserID:     testid,
			TimeNotify: 360000000000,
		}
		testEvent3 := entities.Event{
			ID:         "22112233-4455-6677-8899-aabbccddeeff",
			Title:      "title",
			DateTime:   testdt.AddDate(0, 1, 5),
			Duration:   360000000000,
			Text:       "text",
			UserID:     testid,
			TimeNotify: 360000000000,
		}
		outsideAddMapRepo(repo.m, testEvent)
		outsideAddMapRepo(repo.m, testEvent2)
		outsideAddMapRepo(repo.m, testEvent3)

		var expectedEvents []*entities.Event

		expectedEvents = append(expectedEvents, &testEvent)
		expectedEvents = append(expectedEvents, &testEvent2)

		actualEvents, err := repo.GetForPeriod(ctx, testid, testdt, testdt.AddDate(0, 1, 0))

		require.Nil(t, err)
		require.EqualValues(t, expectedEvents, actualEvents, "events array not expected")
	})
	t.Run("get by period bad", func(t *testing.T) {
		repo.m.Clear()
		testEvent2 := entities.Event{
			ID:         "11112233-4455-6677-8899-aabbccddeeff",
			Title:      "title",
			DateTime:   testdt.AddDate(0, 0, 5),
			Duration:   360000000000,
			Text:       "text",
			UserID:     testid,
			TimeNotify: 360000000000,
		}
		testEvent3 := entities.Event{
			ID:         "22112233-4455-6677-8899-aabbccddeeff",
			Title:      "title",
			DateTime:   testdt.AddDate(0, 1, 5),
			Duration:   360000000000,
			Text:       "text",
			UserID:     testid,
			TimeNotify: 360000000000,
		}
		outsideAddMapRepo(repo.m, testEvent)
		outsideAddMapRepo(repo.m, testEvent2)
		outsideAddMapRepo(repo.m, testEvent3)

		actualEvents, err := repo.GetForPeriod(ctx, testid, testdt.AddDate(0, 1, 0), testdt.AddDate(0, 1, 1))

		require.Nil(t, actualEvents)
		require.Truef(t, errors.Is(err, entities.ErrEventNotFound), "return error must be: %q", entities.ErrEventNotFound)
	})
	t.Run("delete good", func(t *testing.T) {
		repo.m.Clear()
		outsideAddMapRepo(repo.m, testEvent)

		err := repo.DeleteByID(ctx, testid, testid)

		require.Nil(t, err)
	})
	t.Run("delete bad", func(t *testing.T) {
		repo.m.Clear()
		outsideAddMapRepo(repo.m, testEvent)
		err := repo.DeleteByID(ctx, testid, "not id")

		require.Truef(t, errors.Is(err, entities.ErrEventNotFound), "return error must be: %q", entities.ErrEventNotFound)
	})
	t.Run("update good", func(t *testing.T) {
		repo.m.Clear()
		outsideAddMapRepo(repo.m, testEvent)
		updateEvent := entities.Event{
			Title:      "new title",
			DateTime:   testdt,
			Duration:   360000000000,
			Text:       "text",
			UserID:     testid,
			TimeNotify: 360000000000,
		}

		err := repo.UpdateByID(ctx, testid, testid, updateEvent)

		require.Nil(t, err)
		require.Equal(t, "new title", repo.m.events[testid].Title, "updated title not equal real title")
	})
	t.Run("update bad", func(t *testing.T) {
		repo.m.Clear()
		outsideAddMapRepo(repo.m, testEvent)
		err := repo.UpdateByID(ctx, testid, "not id", testEvent)

		require.Truef(t, errors.Is(err, entities.ErrEventNotFound), "return error must be: %q", entities.ErrEventNotFound)
	})
	t.Run("concurrently using", func(t *testing.T) {
		repo.m.Clear()

		wg := &sync.WaitGroup{}
		wg.Add(1)
		echan := make(chan string)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				e := newFakeEvent(i)
				id, err := repo.Add(ctx, e)
				if err == nil {
					echan <- id
				}
			}
			close(echan)
		}()
		deleteCount := 10
		wg.Add(deleteCount)
		for y := 0; y < deleteCount; y++ {
			go func() {
				defer wg.Done()
				for id := range echan {
					repo.DeleteByID(ctx, testid, id)
				}
			}()
		}
		wg.Wait()
		l, _ := repo.GetForPeriod(ctx, testid, time.Now().AddDate(-10, 0, 0), time.Now().AddDate(10, 0, 0))
		require.Equal(t, 0, len(l), "events list must be empty")
	})
}

func newFakeEvent(d int) entities.Event {
	e, _ := entities.NewEvent("title", time.Now().AddDate(0, 0, d).Format(entities.LayoutISO), "6m", "text", testid, "6m")
	return *e
}

func outsideAddMapRepo(m *MapRepo, e entities.Event) {
	m.rwmux.Lock()
	defer m.rwmux.Unlock()
	m.events[e.ID] = &e
	if _, ok := m.users[e.UserID]; !ok {
		m.users[e.UserID] = make(map[time.Time]*entities.Event)
	}
	m.users[e.UserID][e.DateTime] = &e
}
