// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package transport

import (
	"context"
	"github.com/ivch/dynasty/server/handlers/dictionaries"
	"sync"
)

// Ensure, that DictionaryServiceMock does implement DictionaryService.
// If this is not the case, regenerate this file with moq.
var _ DictionaryService = &DictionaryServiceMock{}

// DictionaryServiceMock is a mock implementation of DictionaryService.
//
//	func TestSomethingThatUsesDictionaryService(t *testing.T) {
//
//		// make and configure a mocked DictionaryService
//		mockedDictionaryService := &DictionaryServiceMock{
//			BuildingsListFunc: func(ctx context.Context) ([]*dictionaries.Building, error) {
//				panic("mock out the BuildingsList method")
//			},
//			EntriesListFunc: func(ctx context.Context, buildingID uint) ([]*dictionaries.Entry, error) {
//				panic("mock out the EntriesList method")
//			},
//		}
//
//		// use mockedDictionaryService in code that requires DictionaryService
//		// and then make assertions.
//
//	}
type DictionaryServiceMock struct {
	// BuildingsListFunc mocks the BuildingsList method.
	BuildingsListFunc func(ctx context.Context) ([]*dictionaries.Building, error)

	// EntriesListFunc mocks the EntriesList method.
	EntriesListFunc func(ctx context.Context, buildingID uint) ([]*dictionaries.Entry, error)

	// calls tracks calls to the methods.
	calls struct {
		// BuildingsList holds details about calls to the BuildingsList method.
		BuildingsList []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// EntriesList holds details about calls to the EntriesList method.
		EntriesList []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// BuildingID is the buildingID argument value.
			BuildingID uint
		}
	}
	lockBuildingsList sync.RWMutex
	lockEntriesList   sync.RWMutex
}

// BuildingsList calls BuildingsListFunc.
func (mock *DictionaryServiceMock) BuildingsList(ctx context.Context) ([]*dictionaries.Building, error) {
	if mock.BuildingsListFunc == nil {
		panic("DictionaryServiceMock.BuildingsListFunc: method is nil but DictionaryService.BuildingsList was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockBuildingsList.Lock()
	mock.calls.BuildingsList = append(mock.calls.BuildingsList, callInfo)
	mock.lockBuildingsList.Unlock()
	return mock.BuildingsListFunc(ctx)
}

// BuildingsListCalls gets all the calls that were made to BuildingsList.
// Check the length with:
//
//	len(mockedDictionaryService.BuildingsListCalls())
func (mock *DictionaryServiceMock) BuildingsListCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockBuildingsList.RLock()
	calls = mock.calls.BuildingsList
	mock.lockBuildingsList.RUnlock()
	return calls
}

// EntriesList calls EntriesListFunc.
func (mock *DictionaryServiceMock) EntriesList(ctx context.Context, buildingID uint) ([]*dictionaries.Entry, error) {
	if mock.EntriesListFunc == nil {
		panic("DictionaryServiceMock.EntriesListFunc: method is nil but DictionaryService.EntriesList was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		BuildingID uint
	}{
		Ctx:        ctx,
		BuildingID: buildingID,
	}
	mock.lockEntriesList.Lock()
	mock.calls.EntriesList = append(mock.calls.EntriesList, callInfo)
	mock.lockEntriesList.Unlock()
	return mock.EntriesListFunc(ctx, buildingID)
}

// EntriesListCalls gets all the calls that were made to EntriesList.
// Check the length with:
//
//	len(mockedDictionaryService.EntriesListCalls())
func (mock *DictionaryServiceMock) EntriesListCalls() []struct {
	Ctx        context.Context
	BuildingID uint
} {
	var calls []struct {
		Ctx        context.Context
		BuildingID uint
	}
	mock.lockEntriesList.RLock()
	calls = mock.calls.EntriesList
	mock.lockEntriesList.RUnlock()
	return calls
}
