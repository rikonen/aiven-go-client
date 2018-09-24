package aiven

import (
	"encoding/json"
	"errors"
)

type (
	// Service represents the Service model on Aiven.
	Service struct {
		Backups        []*Backup              `json:"backups"`
		CloudName      string                 `json:"cloud_name"`
		CreateTime     string                 `json:"create_time"`
		UpdateTime     string                 `json:"update_time"`
		GroupList      []string               `json:"group_list"`
		NodeCount      int                    `json:"node_count"`
		Plan           string                 `json:"plan"`
		Name           string                 `json:"service_name"`
		Type           string                 `json:"service_type"`
		URI            string                 `json:"service_uri"`
		URIParams      map[string]string      `json:"service_uri_params"`
		State          string                 `json:"state"`
		Metadata       interface{}            `json:"metadata"`
		Users          []*ServiceUser         `json:"users"`
		UserConfig     map[string]interface{} `json:"user_config"`
		ConnectionInfo ConnectionInfo         `json:"connection_info"`
	}

	// Backup represents an individual backup of service data on Aiven
	Backup struct {
		BackupTime string `json:"backup_time"`
		DataSize   int    `json:"data_size"`
	}

	// ConnectionInfo represents the Service Connection information on Aiven.
	ConnectionInfo struct {
		CassandraHosts []string `json:"cassandra"`

		ElasticsearchURIs     []string `json:"elasticsearch"`
		ElasticsearchUsername string   `json:"elasticsearch_username"`
		ElasticsearchPassword string   `json:"elasticsearch_password"`
		KibanaURI             string   `json:"kibana_uri"`

		GrafanaURIs []string `json:"grafana"`

		InfluxDBURIs         []string `json:"influxdb"`
		InfluxDBDatabaseName string   `json:"influxdb_dbname"`
		InfluxDBUsername     string   `json:"influxdb_username"`
		InfluxDBPassword     string   `json:"influxdb_password"`

		KafkaHosts        []string `json:"kafka"`
		KafkaAccessCert   string   `json:"kafka_access_cert"`
		KafkaAccessKey    string   `json:"kafka_access_key"`
		KafkaConnectURI   string   `json:"kafka_connect_uri"`
		KafkaRestURI      string   `json:"kafka_rest_uri"`
		SchemaRegistryURI string   `json:"schema_registry_uri"`

		PostgresParams      []PostgresParams `json:"pg_params"`
		PostgresReplicaURI  string           `json:"pg_replica_uri"`
		PostgresStandbyURIs []string         `json:"pg_standby"`
		PostgresURIs        []string         `json:"pg"`

		RedisPassword  string   `json:"redis_password"`
		RedisSlaveURIs []string `json:"redis_slave"`
		RedisURIs      []string `json:"redis"`
	}

	// PostgresParams represents individual parameters for a PostgreSQL ConnectionInfo
	PostgresParams struct {
		DatabaseName string `json:"dbname"`
		Host         string `string:"host"`
		Password     string `string:"password"`
		Port         string `string:"port"`
		SSLMode      string `string:"sslmode"`
		User         string `string:"user"`
	}

	// ServicesHandler is the client that interacts with the Service API
	// endpoints on Aiven.
	ServicesHandler struct {
		client *Client
	}

	// CreateServiceRequest are the parameters to create a Service.
	CreateServiceRequest struct {
		Cloud       string                 `json:"cloud,omitempty"`
		GroupName   string                 `json:"group_name,omitempty"`
		Plan        string                 `json:"plan,omitempty"`
		ServiceName string                 `json:"service_name"`
		ServiceType string                 `json:"service_type"`
		UserConfig  map[string]interface{} `json:"user_config,omitempty"`
	}

	// UpdateServiceRequest are the parameters to update a Service.
	UpdateServiceRequest struct {
		Cloud      string                 `json:"cloud,omitempty"`
		GroupName  string                 `json:"group_name,omitempty"`
		Plan       string                 `json:"plan,omitempty"`
		Powered    bool                   `json:"powered"` // TODO: figure out if we can overwrite the default?
		UserConfig map[string]interface{} `json:"user_config,omitempty"`
	}

	// ServiceResponse represents the response from Aiven after interacting with
	// the Service API.
	ServiceResponse struct {
		APIResponse
		Service *Service `json:"service"`
	}

	// ServiceListResponse represents the response from Aiven for listing
	// services.
	ServiceListResponse struct {
		APIResponse
		Services []*Service `json:"services"`
	}
)

// Hostname provides host name for the service. This method is provided for backwards
// compatibility, typically it is easier to just get the value from URIParams directly.
func (s *Service) Hostname() (string, error) {
	return s.URIParams["host"], nil
}

// Port provides port for the service. This method is provided for backwards
// compatibility, typically it is easier to just get the value from URIParams directly.
func (s *Service) Port() (string, error) {
	return s.URIParams["port"], nil
}

// Create creates the given Service on Aiven.
func (h *ServicesHandler) Create(project string, req CreateServiceRequest) (*Service, error) {
	path := buildPath("project", project, "service")
	rsp, err := h.client.doPostRequest(path, req)
	if err != nil {
		return nil, err
	}

	return parseServiceResponse(rsp)
}

// Get gets a specific service from Aiven.
func (h *ServicesHandler) Get(project, service string) (*Service, error) {
	path := buildPath("project", project, "service", service)
	rsp, err := h.client.doGetRequest(path, nil)
	if err != nil {
		return nil, err
	}

	return parseServiceResponse(rsp)
}

// Update will update the given service with the given parameters.
func (h *ServicesHandler) Update(project, service string, req UpdateServiceRequest) (*Service, error) {
	path := buildPath("project", project, "service", service)
	rsp, err := h.client.doPutRequest(path, req)
	if err != nil {
		return nil, err
	}

	return parseServiceResponse(rsp)
}

// Delete will delete the given service from Aiven.
func (h *ServicesHandler) Delete(project, service string) error {
	path := buildPath("project", project, "service", service)
	bts, err := h.client.doDeleteRequest(path, nil)
	if err != nil {
		return err
	}

	return handleDeleteResponse(bts)
}

// List will fetch all services for a given project.
func (h *ServicesHandler) List(project string) ([]*Service, error) {
	path := buildPath("project", project, "service")
	rsp, err := h.client.doGetRequest(path, nil)
	if err != nil {
		return nil, err
	}

	var response *ServiceListResponse
	if err := json.Unmarshal(rsp, &response); err != nil {
		return nil, err
	}

	if len(response.Errors) != 0 {
		return nil, errors.New(response.Message)
	}

	return response.Services, nil
}

func parseServiceResponse(rsp []byte) (*Service, error) {
	var response *ServiceResponse
	if err := json.Unmarshal(rsp, &response); err != nil {
		return nil, err
	}

	if len(response.Errors) != 0 {
		return nil, errors.New(response.Message)
	}

	return response.Service, nil
}
