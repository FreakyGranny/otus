package memorystorage

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		s := New()
		e := &storage.Event{
			Title: "test",
		}
		err := s.CreateEvent(e)
		require.NoError(t, err)
		require.NotEqual(t, 0, e.ID)
	})
	t.Run("get", func(t *testing.T) {
		s := New()
		e := &storage.Event{
			Title: "test",
		}
		err := s.CreateEvent(e)
		require.NoError(t, err)
		require.NotEqual(t, 0, e.ID)

		res, err := s.GetEvent(e.ID)
		require.NoError(t, err)
		require.Equal(t, e.ID, res.ID)
	})
	t.Run("get not found", func(t *testing.T) {
		s := New()

		res, err := s.GetEvent(65635)
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("delete", func(t *testing.T) {
		s := New()
		e := &storage.Event{
			Title: "test",
		}
		err := s.CreateEvent(e)
		require.NoError(t, err)
		require.NotEqual(t, 0, e.ID)

		err = s.DeleteEvent(e.ID)
		require.NoError(t, err)
	})
	t.Run("update", func(t *testing.T) {
		s := New()
		e := &storage.Event{
			Title: "test",
		}
		err := s.CreateEvent(e)
		require.NoError(t, err)
		require.NotEqual(t, 0, e.ID)

		nEvent := &storage.Event{
			ID: e.ID,
			Title: "updated",
		}
		err = s.UpdateEvent(nEvent)
		require.NoError(t, err)

		res, err := s.GetEvent(nEvent.ID)
		require.NoError(t, err)
		require.Equal(t, nEvent.Title, res.Title)
	})
	t.Run("get list", func(t *testing.T) {
		s := New()
		fe := &storage.Event{
			Title: "test1",
		}
		se := &storage.Event{
			Title: "test2",
		}

		err := s.CreateEvent(fe)
		require.NoError(t, err)
		err = s.CreateEvent(se)
		require.NoError(t, err)

		res, err := s.GetEventList()
		require.NoError(t, err)
		require.Equal(t, 2, len(res))
		require.Equal(t, res, []*storage.Event{fe, se})
	})
}
