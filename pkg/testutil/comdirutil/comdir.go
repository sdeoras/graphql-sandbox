package comdirutil

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/sdeoras/graphql/pkg/log"
	"github.com/sdeoras/graphql/pkg/rest/mw/auth"
	"github.com/sdeoras/graphql/pkg/rest/mw/auth/authenticator"
	"github.com/sdeoras/graphql/pkg/rest/mw/auth/authorizer"
	"math"
	"time"
)

const (
	query      = "query"
	count      = "count"
	id         = "id"
	name       = "name"
	joinDate   = "joinDate"
	endDate    = "endDate"
	employee   = "employee"
	employees  = "employees"
	manager    = "maanger"
	manages    = "manages"
	contractor = "contractor"
)

type Request struct {
	ID string `json:"id"`
}

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Meta struct {
	JoinDate   time.Time `json:"joinDate"`
	EndDate    time.Time `json:"endDate"`
	Department string    `json:"department"`
}

type Employee struct {
	*Person
	*Meta
	Manager string   `json:"manager"`
	Manages []string `json:"manages"`
}

type Contractor struct {
	*Person
	*Meta
}

var Schema graphql.Schema

var (
	departmentEnum  *graphql.Enum
	employeeType    *graphql.Object
	contractorType  *graphql.Object
	personInterface *graphql.Interface
)

func init() {
	authN := authenticator.NewAuthenticator(&authenticator.Config{
		AllowedUsers:  []string{},
		AllowedGroups: []string{auth.GroupGoogle},
		Logger:        log.Logger(),
	})
	authZ := authorizer.NewAuthorizer(&authorizer.Config{
		AllowedRoles: []string{auth.RoleAdmin},
		Logger:       log.Logger(),
	})

	request := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "request",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Name:              "id",
				Type:              graphql.String,
				Args:              nil,
				Resolve:           nil,
				DeprecationReason: "",
				Description:       "",
			},
		},
		Description: "id to query something with it",
	})

	_ = request

	employeeType = graphql.NewObject(graphql.ObjectConfig{
		Name: employee,
		// employee object refers to list of employees (i.e., it references self)
		// this is a special case:
		// https://github.com/graphql-go/graphql/issues/112#issuecomment-680168865
		Fields: (graphql.FieldsThunk)(func() graphql.Fields {
			return graphql.Fields{
				id: &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
					Args: nil,
					Resolve: authN.Authenticate(
						authZ.Authorize(
							func(p graphql.ResolveParams) (interface{}, error) {
								switch v := p.Source.(type) {
								case *Employee:
									return v.ID, nil
								default:
									return nil, fmt.Errorf("invalid object type %T, expected *Employee", v)
								}
							},
						),
					),
					DeprecationReason: "",
					Description:       "id of the employee",
				},
				name: &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
					Args: nil,
					Resolve: authN.Authenticate(
						authZ.Authorize(
							func(p graphql.ResolveParams) (interface{}, error) {
								switch v := p.Source.(type) {
								case *Employee:
									return v.Name, nil
								default:
									return nil, fmt.Errorf("invalid object type %T, expected *Employee", v)
								}
							},
						),
					),
					DeprecationReason: "",
					Description:       "name of the employee",
				},
				joinDate: &graphql.Field{
					Name: joinDate,
					Type: graphql.String,
					Args: nil,
					Resolve: authN.Authenticate(
						authZ.Authorize(
							func(p graphql.ResolveParams) (interface{}, error) {
								switch v := p.Source.(type) {
								case *Employee:
									return v.JoinDate, nil
								default:
									return nil, fmt.Errorf("invalid object type %T, expected *Employee", v)
								}
							},
						),
					),
					DeprecationReason: "",
					Description:       "join date of the employee",
				},
				endDate: &graphql.Field{
					Name: endDate,
					Type: graphql.String,
					Args: nil,
					Resolve: authN.Authenticate(
						authZ.Authorize(
							func(p graphql.ResolveParams) (interface{}, error) {
								switch v := p.Source.(type) {
								case *Employee:
									return v.EndDate, nil
								default:
									return nil, fmt.Errorf("invalid object type %T, expected *Employee", v)
								}
							},
						),
					),
					DeprecationReason: "",
					Description:       "termination date of the employee",
				},
				manager: &graphql.Field{
					Name: manager,
					Type: employeeType,
					Args: nil,
					Resolve: authN.Authenticate(
						authZ.Authorize(
							func(p graphql.ResolveParams) (interface{}, error) {
								switch v := p.Source.(type) {
								case *Employee:
									e, ok := registry[v.Manager]
									if !ok {
										return nil, fmt.Errorf("manager id not found")
									}
									return e, nil
								default:
									return nil, fmt.Errorf("invalid object type %T, expected *Employee", v)
								}
							},
						),
					),
					DeprecationReason: "",
					Description:       "manager of the employee",
				},
				manages: &graphql.Field{
					Name: manages,
					Type: graphql.NewList(employeeType),
					Args: nil,
					Resolve: authN.Authenticate(
						authZ.Authorize(
							func(p graphql.ResolveParams) (interface{}, error) {
								switch v := p.Source.(type) {
								case *Employee:
									var employees []*Employee
									for _, eid := range v.Manages {
										e, ok := registry[eid]
										if !ok {
											return nil, fmt.Errorf("manager id not found")
										}
										employees = append(employees, e)
									}
									return employees, nil
								default:
									return nil, fmt.Errorf("invalid object type %T, expected *Employee", v)
								}
							},
						),
					),
					DeprecationReason: "",
					Description:       "list of employees this employee manages",
				},
			}
		}),
		IsTypeOf:    nil,
		Description: "type def for employee",
	})

	contractorType = graphql.NewObject(graphql.ObjectConfig{
		Name: "contractor type",
		Interfaces: []*graphql.Interface{
			personInterface,
		},
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Name: "id of the contractor",
				Type: graphql.String,
				Args: nil,
				Resolve: authN.Authenticate(
					authZ.Authorize(
						func(p graphql.ResolveParams) (interface{}, error) {
							switch v := p.Source.(type) {
							case *Contractor:
								return v.ID, nil
							default:
								return nil, fmt.Errorf("invalid object type %T, expected *Contractor", v)
							}
						},
					),
				),
				DeprecationReason: "",
				Description:       "",
			},
			"name": &graphql.Field{
				Name: "name of the contractor",
				Type: graphql.String,
				Args: nil,
				Resolve: authN.Authenticate(
					authZ.Authorize(
						func(p graphql.ResolveParams) (interface{}, error) {
							switch v := p.Source.(type) {
							case *Contractor:
								return v.Name, nil
							default:
								return nil, fmt.Errorf("invalid object type %T, expected *Contractor", v)
							}
						},
					),
				),
				DeprecationReason: "",
				Description:       "",
			},
			"joinDate": &graphql.Field{
				Name: "joining date of the contractor",
				Type: graphql.String,
				Args: nil,
				Resolve: authN.Authenticate(
					authZ.Authorize(
						func(p graphql.ResolveParams) (interface{}, error) {
							switch v := p.Source.(type) {
							case *Contractor:
								return v.JoinDate, nil
							default:
								return nil, fmt.Errorf("invalid object type %T, expected *Contractor", v)
							}
						},
					),
				),
				DeprecationReason: "",
				Description:       "",
			},
			"endDate": &graphql.Field{
				Name: "end date of the contractor",
				Type: graphql.String,
				Args: nil,
				Resolve: authN.Authenticate(
					authZ.Authorize(
						func(p graphql.ResolveParams) (interface{}, error) {
							switch v := p.Source.(type) {
							case *Contractor:
								return v.EndDate, nil
							default:
								return nil, fmt.Errorf("invalid object type %T, expected *Contractor", v)
							}
						},
					),
				),
				DeprecationReason: "",
				Description:       "",
			},
			"department": &graphql.Field{
				Name: "department of the contractor",
				Type: graphql.String,
				Args: nil,
				Resolve: authN.Authenticate(
					authZ.Authorize(
						func(p graphql.ResolveParams) (interface{}, error) {
							switch v := p.Source.(type) {
							case *Contractor:
								return v.Department, nil
							default:
								return nil, fmt.Errorf("invalid object type %T, expected *Contractor", v)
							}
						},
					),
				),
				DeprecationReason: "",
				Description:       "",
			},
		},
		IsTypeOf:    nil,
		Description: "type def for contractor",
	})

	departmentEnum = graphql.NewEnum(graphql.EnumConfig{
		Name: "Department",
		Values: graphql.EnumValueConfigMap{
			"Engineering": &graphql.EnumValueConfig{
				Value:             "eng",
				DeprecationReason: "",
				Description:       "engineering department",
			},
			"Finance": &graphql.EnumValueConfig{
				Value:             "fin",
				DeprecationReason: "",
				Description:       "finance department",
			},
			"HumanResources": &graphql.EnumValueConfig{
				Value:             "hr",
				DeprecationReason: "",
				Description:       "HR department",
			},
		},
		Description: "one of the departments the employee belongs to",
	})

	personInterface = graphql.NewInterface(
		graphql.InterfaceConfig{
			Name: "someone who is a person",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Name: "id of the entity",
					Type: graphql.String,
					Args: nil,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						switch v := p.Source.(type) {
						case Person:
							return v.ID, nil
						default:
							return nil, fmt.Errorf("invalid object type %T, expected Person", v)
						}
					},
					DeprecationReason: "",
					Description:       "id of the person",
				},
				"name": &graphql.Field{
					Name: "name of the entity",
					Type: graphql.String,
					Args: nil,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						switch v := p.Source.(type) {
						case Person:
							return v.Name, nil
						default:
							return nil, fmt.Errorf("invalid object type %T, expected Person", v)
						}
					},
					DeprecationReason: "",
					Description:       "name of the person",
				},
			},
			ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
				switch p.Value.(type) {
				case *Employee:
					return employeeType
				case *Contractor:
					return contractorType
				default:
					return nil
				}
			},
			Description: "typedef for person interface",
		},
	)

	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: query,
			Fields: graphql.Fields{
				employee: &graphql.Field{
					Name: employee,
					Type: employeeType,
					Args: graphql.FieldConfigArgument{
						id: &graphql.ArgumentConfig{
							Type:         graphql.String,
							DefaultValue: nil,
							Description:  "id of the employee to fetch",
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						eid, ok := p.Args[id].(string)
						if !ok {
							return nil, fmt.Errorf("invalid id type")
						}
						e, ok := registry[eid]
						if !ok {
							return nil, fmt.Errorf("manager id not found")
						}

						return e, nil
					},
					DeprecationReason: "",
					Description:       "fetch employee info based on id",
				},
				employees: &graphql.Field{
					Name: employees,
					Type: graphql.NewList(employeeType),
					Args: graphql.FieldConfigArgument{
						count: &graphql.ArgumentConfig{
							Type:         graphql.Int,
							DefaultValue: nil,
							Description:  "count how many (up to)",
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						n := math.MaxInt32
						if m, ok := p.Args[count].(int); ok && m < n {
							n = m
						}

						var e []*Employee
						for _, v := range registry {
							if len(e) >= n {
								break
							}
							v := v
							e = append(e, v)
						}

						return e, nil
					},
					DeprecationReason: "",
					Description:       "list all employees",
				},
			},
		},
	)
	Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
