export PORT="9001"
export services="biggraph,twowaykv,links,geoip,wikipedia,analytics,grafana,prometheus"

###################
## core services ##
###################

export biggraph_incoming_path="/services/biggraph/"
export biggraph_outgoing_url="http://localhost:5000"

export twowaykv_incoming_path="/services/twowaykv/"
export twowaykv_outgoing_url="http://localhost:5001"

export wikipedia_incoming_path="/services/wiki/"
export wikipedia_outgoing_url="https://en.wikipedia.org"

export links_incoming_path="/"
export links_outgoing_url="http://localhost:3000"

###############
## analytics ##
###############

export geoip_incoming_path="/analytics/api/geoIpServer/"
export geoip_outgoing_url="http://api.ipstack.com"

export analytics_incoming_path="/analytics/server/"
export analytics_outgoing_url="http://localhost:5002"

###########
## admin ##
###########

export grafana_incoming_path="/admin/grafana/"
export grafana_outgoing_url="http://localhost:5000"

export prometheus_incoming_path="/admin/prometheus/"
export prometheus_outgoing_url="http://localhost:9090"
