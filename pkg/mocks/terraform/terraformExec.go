// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"

	tfexec "github.com/hashicorp/terraform-exec/tfexec"

	tfjson "github.com/hashicorp/terraform-json"
)

// TerraformExec is an autogenerated mock type for the terraformExec type
type TerraformExec struct {
	mock.Mock
}

// Apply provides a mock function with given fields: ctx, opts
func (_m *TerraformExec) Apply(ctx context.Context, opts ...tfexec.ApplyOption) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...tfexec.ApplyOption) error); ok {
		r0 = rf(ctx, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Init provides a mock function with given fields: ctx, opts
func (_m *TerraformExec) Init(ctx context.Context, opts ...tfexec.InitOption) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...tfexec.InitOption) error); ok {
		r0 = rf(ctx, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Plan provides a mock function with given fields: ctx, opts
func (_m *TerraformExec) Plan(ctx context.Context, opts ...tfexec.PlanOption) (bool, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, ...tfexec.PlanOption) bool); ok {
		r0 = rf(ctx, opts...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...tfexec.PlanOption) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetEnv provides a mock function with given fields: env
func (_m *TerraformExec) SetEnv(env map[string]string) error {
	ret := _m.Called(env)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]string) error); ok {
		r0 = rf(env)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetStdout provides a mock function with given fields: w
func (_m *TerraformExec) SetStdout(w io.Writer) {
	_m.Called(w)
}

// Validate provides a mock function with given fields: ctx
func (_m *TerraformExec) Validate(ctx context.Context) (*tfjson.ValidateOutput, error) {
	ret := _m.Called(ctx)

	var r0 *tfjson.ValidateOutput
	if rf, ok := ret.Get(0).(func(context.Context) *tfjson.ValidateOutput); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tfjson.ValidateOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WorkspaceNew provides a mock function with given fields: ctx, workspace, opts
func (_m *TerraformExec) WorkspaceNew(ctx context.Context, workspace string, opts ...tfexec.WorkspaceNewCmdOption) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, workspace)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...tfexec.WorkspaceNewCmdOption) error); ok {
		r0 = rf(ctx, workspace, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WorkspaceSelect provides a mock function with given fields: ctx, workspace
func (_m *TerraformExec) WorkspaceSelect(ctx context.Context, workspace string) error {
	ret := _m.Called(ctx, workspace)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, workspace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTerraformExec interface {
	mock.TestingT
	Cleanup(func())
}

// NewTerraformExec creates a new instance of TerraformExec. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTerraformExec(t mockConstructorTestingTNewTerraformExec) *TerraformExec {
	mock := &TerraformExec{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
