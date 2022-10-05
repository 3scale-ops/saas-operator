{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "editable": true,
  "gnetId": null,
  "graphTooltip": 0,
  "id": 1,
  "links": [],
  "panels": [
    {
      "collapsed": false,
      "datasource": null,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 0
      },
      "id": 20,
      "panels": [],
      "repeat": "deployment",
      "scopedVars": {
        "deployment": {}
      },
      "title": "Twemproxy $deployment",
      "type": "row"
    },
    {
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "thresholds"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 1
      },
      "id": 47,
      "links": [],
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "center",
        "orientation": "auto",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "text": {},
        "textMode": "value"
      },
      "pluginVersion": "7.5.15",
      "scopedVars": {
        "deployment": {}
      },
      "targets": [
        {
          "exemplar": true,
          "expr": "count(twemproxy_exporter_up{namespace=\"3scale-saas\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\"} == 1) or on() vector(0)",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "timeFrom": null,
      "timeShift": null,
      "title": "UP",
      "type": "stat"
    },
    {
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "thresholds"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 1
      },
      "id": 66,
      "links": [],
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "center",
        "orientation": "auto",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "text": {},
        "textMode": "value"
      },
      "pluginVersion": "7.5.15",
      "scopedVars": {
        "deployment": {}
      },
      "targets": [
        {
          "exemplar": true,
          "expr": "count(twemproxy_exporter_up{namespace=\"3scale-saas\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\"} == 0) or on() vector(0)",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "",
          "hide": false,
          "interval": "",
          "legendFormat": "",
          "refId": "B"
        }
      ],
      "timeFrom": null,
      "timeShift": null,
      "title": "DOWN",
      "type": "stat"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 9
      },
      "hiddenSeries": false,
      "id": 39,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_requests_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_requests_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        },
        {
          "hide": false,
          "refId": "C"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Requests Total  requests/second (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 9
      },
      "hiddenSeries": false,
      "id": 65,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_requests_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_requests_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Requests Total  requests/second (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 17
      },
      "hiddenSeries": false,
      "id": 91,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_responses_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_responses_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        },
        {
          "hide": false,
          "refId": "C"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Responses Total  requests/second (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 17
      },
      "hiddenSeries": false,
      "id": 92,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_responses_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_responses_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        },
        {
          "hide": false,
          "refId": "C"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Responses Total  requests/second (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 25
      },
      "hiddenSeries": false,
      "id": 93,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_requests_bytes_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_requests_bytes_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        },
        {
          "hide": false,
          "refId": "C"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Requests Total bytes (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 25
      },
      "hiddenSeries": false,
      "id": 94,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_requests_bytes_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_requests_bytes_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        },
        {
          "hide": false,
          "refId": "C"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Requests Total bytes (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 33
      },
      "hiddenSeries": false,
      "id": 95,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_responses_bytes_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_responses_bytes_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        },
        {
          "hide": false,
          "refId": "C"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Responses Total bytes (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 33
      },
      "hiddenSeries": false,
      "id": 96,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_responses_bytes_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_responses_bytes_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        },
        {
          "hide": false,
          "refId": "C"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Responses Total bytes (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 41
      },
      "hiddenSeries": false,
      "id": 67,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "twemproxy_exporter_current_connections{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\"}",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_current_connections{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Current Connections (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 41
      },
      "hiddenSeries": false,
      "id": 89,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "twemproxy_exporter_client_connections_active{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_client_connections_active{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Client Connections Active (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 49
      },
      "hiddenSeries": false,
      "id": 68,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_server_connections_active{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_server_connections_active{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Connections Active (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 49
      },
      "hiddenSeries": false,
      "id": 90,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_server_connections_active{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_server_connections_active{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Connections Active (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 57
      },
      "hiddenSeries": false,
      "id": 69,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "rate(twemproxy_exporter_client_err_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_client_err_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Client Errors Total request/second (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 57
      },
      "hiddenSeries": false,
      "id": 70,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "rate(twemproxy_exporter_client_eof_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_client_eof_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Client Eof Total request/second (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 65
      },
      "hiddenSeries": false,
      "id": 71,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_err_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_err_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Error Total requests/second (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 65
      },
      "hiddenSeries": false,
      "id": 72,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_err_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_err_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Error Total requests/second (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 73
      },
      "hiddenSeries": false,
      "id": 73,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_eof_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_eof_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server EoF Total requests/second (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 73
      },
      "hiddenSeries": false,
      "id": 74,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_eof_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_eof_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server EoF Total requests/second (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 81
      },
      "hiddenSeries": false,
      "id": 75,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_timeouts_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_timeouts_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Timeouts Total requests/second (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 81
      },
      "hiddenSeries": false,
      "id": 76,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_timeouts_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_server_timeouts_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Server Timeouts Total requests/second (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 89
      },
      "hiddenSeries": false,
      "id": 77,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "rate(twemproxy_exporter_backend_server_ejections_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_backend_server_ejections_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Backend Server Ejections (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 89
      },
      "hiddenSeries": false,
      "id": 78,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "rate(twemproxy_exporter_forward_errors_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m])",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(rate(twemproxy_exporter_forward_errors_total{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}[1m]))",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Forward Error Total (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 97
      },
      "hiddenSeries": false,
      "id": 79,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_incoming_queue{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_incoming_queue{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Incoming queue requests/second (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 97
      },
      "hiddenSeries": false,
      "id": 80,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_incoming_queue{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_incoming_queue{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Incoming queue requests/second (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 105
      },
      "hiddenSeries": false,
      "id": 81,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_outgoing_queue{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_outgoing_queue{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Outgoing queue requests/second (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 105
      },
      "hiddenSeries": false,
      "id": 82,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_outgoing_queue{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_outgoing_queue{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Outgoing queue requests/second (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "description": "",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 113
      },
      "hiddenSeries": false,
      "id": 83,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_incoming_queue_bytes{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_incoming_queue_bytes{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Incoming queue bytes (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "description": "",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 113
      },
      "hiddenSeries": false,
      "id": 84,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_incoming_queue_bytes{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_incoming_queue_bytes{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Incoming queue bytes (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "description": "",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 121
      },
      "hiddenSeries": false,
      "id": 85,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_outgoing_queue_bytes{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}) by (server)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{server}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_outgoing_queue_bytes{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Outgoing queue bytes (by shard)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "description": "",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 121
      },
      "hiddenSeries": false,
      "id": 86,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [
        {
          "alias": "total",
          "fill": 4,
          "linewidth": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_outgoing_queue_bytes{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"}) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        },
        {
          "exemplar": true,
          "expr": "sum(twemproxy_exporter_outgoing_queue_bytes{namespace=\"$namespace\",pod=~\"[[deployment]]-[a-z0-9]+-[a-z0-9]+\",pool=\"[[pool]]\"})",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "total",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Outgoing queue bytes (by pod)",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "decimals": 0,
          "format": "decbytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "collapsed": false,
      "datasource": null,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 129
      },
      "id": 13,
      "panels": [],
      "repeat": "deployment",
      "scopedVars": {
        "deployment": {}
      },
      "title": "Pods ($deployment)",
      "type": "row"
    },
    {
      "cacheTimeout": null,
      "colorBackground": true,
      "colorValue": false,
      "colors": [
        "#F2495C",
        "rgba(237, 129, 40, 0.89)",
        "#299c46"
      ],
      "datasource": "$datasource",
      "decimals": 0,
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "format": "none",
      "gauge": {
        "maxValue": 100,
        "minValue": 0,
        "show": false,
        "thresholdLabels": false,
        "thresholdMarkers": true
      },
      "gridPos": {
        "h": 3,
        "w": 6,
        "x": 0,
        "y": 130
      },
      "hideTimeOverride": true,
      "id": 30,
      "interval": "",
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "name": "value to text",
          "value": 1
        },
        {
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "pluginVersion": "6.2.4",
      "postfix": "",
      "postfixFontSize": "50%",
      "prefix": "",
      "prefixFontSize": "50%",
      "rangeMaps": [
        {
          "from": "null",
          "text": "N/A",
          "to": "null"
        }
      ],
      "scopedVars": {
        "deployment": {}
      },
      "sparkline": {
        "fillColor": "rgba(31, 118, 189, 0.18)",
        "full": false,
        "lineColor": "rgb(31, 120, 193)",
        "show": false
      },
      "tableColumn": "",
      "targets": [
        {
          "expr": "sum(kube_deployment_status_replicas_available{namespace='$namespace',deployment=~'$deployment'})",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "thresholds": "1,2",
      "timeFrom": "30s",
      "timeShift": "30s",
      "title": "Running pods",
      "type": "singlestat",
      "valueFontSize": "80%",
      "valueMaps": [
        {
          "op": "=",
          "text": "0",
          "value": "null"
        }
      ],
      "valueName": "avg"
    },
    {
      "cacheTimeout": null,
      "colorBackground": true,
      "colorPrefix": false,
      "colorValue": false,
      "colors": [
        "#299c46",
        "rgba(237, 129, 40, 0.89)",
        "#F2495C"
      ],
      "datasource": "$datasource",
      "decimals": 0,
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "format": "none",
      "gauge": {
        "maxValue": 100,
        "minValue": 0,
        "show": false,
        "thresholdLabels": false,
        "thresholdMarkers": true
      },
      "gridPos": {
        "h": 3,
        "w": 6,
        "x": 6,
        "y": 130
      },
      "hideTimeOverride": true,
      "id": 32,
      "interval": null,
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "name": "value to text",
          "value": 1
        },
        {
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "postfix": "",
      "postfixFontSize": "50%",
      "prefix": "",
      "prefixFontSize": "50%",
      "rangeMaps": [
        {
          "from": "null",
          "text": "N/A",
          "to": "null"
        }
      ],
      "scopedVars": {
        "deployment": {}
      },
      "sparkline": {
        "fillColor": "rgba(31, 118, 189, 0.18)",
        "full": false,
        "lineColor": "rgb(31, 120, 193)",
        "show": false
      },
      "tableColumn": "",
      "targets": [
        {
          "expr": "sum(kube_deployment_status_replicas_unavailable{namespace='$namespace',deployment=~'$deployment'})",
          "format": "time_series",
          "intervalFactor": 1,
          "refId": "A"
        }
      ],
      "thresholds": "1,2",
      "timeFrom": "30s",
      "timeShift": "30s",
      "title": "Unavailable pods",
      "type": "singlestat",
      "valueFontSize": "80%",
      "valueMaps": [
        {
          "op": "=",
          "text": "0",
          "value": "null"
        }
      ],
      "valueName": "avg"
    },
    {
      "cacheTimeout": null,
      "colorBackground": true,
      "colorValue": false,
      "colors": [
        "#F2495C",
        "rgba(237, 129, 40, 0.89)",
        "#299c46"
      ],
      "datasource": "$datasource",
      "decimals": 0,
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "format": "none",
      "gauge": {
        "maxValue": 100,
        "minValue": 0,
        "show": false,
        "thresholdLabels": false,
        "thresholdMarkers": true
      },
      "gridPos": {
        "h": 3,
        "w": 6,
        "x": 12,
        "y": 130
      },
      "hideTimeOverride": true,
      "id": 37,
      "interval": "",
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "name": "value to text",
          "value": 1
        },
        {
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "pluginVersion": "6.2.4",
      "postfix": "",
      "postfixFontSize": "50%",
      "prefix": "",
      "prefixFontSize": "50%",
      "rangeMaps": [
        {
          "from": "null",
          "text": "N/A",
          "to": "null"
        }
      ],
      "scopedVars": {
        "deployment": {}
      },
      "sparkline": {
        "fillColor": "rgba(31, 118, 189, 0.18)",
        "full": false,
        "lineColor": "rgb(31, 120, 193)",
        "show": false
      },
      "tableColumn": "",
      "targets": [
        {
          "expr": "count(count(container_memory_working_set_bytes{namespace='$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (node))",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "thresholds": "1,2",
      "timeFrom": "30s",
      "timeShift": "30s",
      "title": "Pods distributed on hosts",
      "type": "singlestat",
      "valueFontSize": "80%",
      "valueMaps": [
        {
          "op": "=",
          "text": "0",
          "value": "null"
        }
      ],
      "valueName": "avg"
    },
    {
      "cacheTimeout": null,
      "colorBackground": true,
      "colorValue": false,
      "colors": [
        "#299c46",
        "rgba(237, 129, 40, 0.89)",
        "#d44a3a"
      ],
      "datasource": "$datasource",
      "decimals": 0,
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "format": "none",
      "gauge": {
        "maxValue": 100,
        "minValue": 0,
        "show": false,
        "thresholdLabels": false,
        "thresholdMarkers": true
      },
      "gridPos": {
        "h": 3,
        "w": 6,
        "x": 18,
        "y": 130
      },
      "hideTimeOverride": true,
      "id": 36,
      "interval": null,
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "name": "value to text",
          "value": 1
        },
        {
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "postfix": "",
      "postfixFontSize": "50%",
      "prefix": "",
      "prefixFontSize": "50%",
      "rangeMaps": [
        {
          "from": "null",
          "text": "N/A",
          "to": "null"
        }
      ],
      "scopedVars": {
        "deployment": {}
      },
      "sparkline": {
        "fillColor": "rgba(31, 118, 189, 0.18)",
        "full": false,
        "lineColor": "rgb(31, 120, 193)",
        "show": false
      },
      "tableColumn": "",
      "targets": [
        {
          "expr": "max(sum(delta(kube_pod_container_status_restarts_total{namespace='$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}[5m])) by (pod))",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "thresholds": "1,2",
      "timeFrom": "30s",
      "timeShift": "30s",
      "title": "Max pods restarts (last 5 minutes)",
      "type": "singlestat",
      "valueFontSize": "80%",
      "valueMaps": [
        {
          "op": "=",
          "text": "0",
          "value": "null"
        }
      ],
      "valueName": "avg"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 7,
        "w": 24,
        "x": 0,
        "y": 133
      },
      "hiddenSeries": false,
      "id": 11,
      "legend": {
        "avg": false,
        "current": false,
        "hideEmpty": true,
        "hideZero": true,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null as zero",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "kube_deployment_status_replicas{namespace='$namespace',deployment=~'$deployment'}",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "{{`{{deployment}}`}}-total-pods",
          "legendLink": null,
          "refId": "A",
          "step": 10
        },
        {
          "expr": "kube_deployment_status_replicas_available{namespace='$namespace',deployment=~'$deployment'}",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "{{`{{deployment}}`}}-avail-pods",
          "refId": "B"
        },
        {
          "expr": "kube_deployment_status_replicas_unavailable{namespace='$namespace',deployment=~'$deployment'}",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "{{`{{deployment}}`}}-unavail-pods",
          "refId": "C"
        },
        {
          "expr": "count(count(container_memory_working_set_bytes{namespace='$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (node))",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "used-hosts",
          "refId": "D"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Pod count (total, avail, unvail) and pods hosts distribution",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": 0,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": false
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 6,
        "w": 24,
        "x": 0,
        "y": 140
      },
      "hiddenSeries": false,
      "id": 9,
      "legend": {
        "avg": false,
        "current": false,
        "hideEmpty": true,
        "hideZero": true,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "repeatedByRow": false,
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(delta(kube_pod_container_status_restarts_total{namespace='$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}[5m])) by (pod)",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Pods restarts (last 5 minutes)",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "collapsed": false,
      "datasource": null,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 146
      },
      "id": 4,
      "panels": [],
      "repeat": "deployment",
      "scopedVars": {
        "deployment": {}
      },
      "title": "CPU Usage ($deployment)",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 7,
        "w": 24,
        "x": 0,
        "y": 147
      },
      "hiddenSeries": false,
      "id": 64,
      "interval": "",
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null as zero",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "legendLink": null,
          "refId": "A",
          "step": 10
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "CPU Usage",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": 0,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": false
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "collapsed": false,
      "datasource": null,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 154
      },
      "id": 5,
      "panels": [],
      "repeat": "deployment",
      "scopedVars": {
        "deployment": {}
      },
      "title": "CPU Quota ($deployment)",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "columns": [],
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fontSize": "100%",
      "gridPos": {
        "h": 7,
        "w": 24,
        "x": 0,
        "y": 155
      },
      "id": 1,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null as zero",
      "pageSize": null,
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "scroll": true,
      "seriesOverrides": [],
      "showHeader": true,
      "sort": {
        "col": 1,
        "desc": false
      },
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "styles": [
        {
          "alias": "Time",
          "align": "auto",
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "pattern": "Time",
          "type": "hidden"
        },
        {
          "alias": "CPU Usage",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": false,
          "linkTooltip": "Drill down",
          "linkUrl": "",
          "pattern": "Value #A",
          "thresholds": [],
          "type": "number",
          "unit": "short"
        },
        {
          "alias": "CPU Requests",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": false,
          "linkTooltip": "Drill down",
          "linkUrl": "",
          "pattern": "Value #B",
          "thresholds": [],
          "type": "number",
          "unit": "short"
        },
        {
          "alias": "CPU Requests %",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": false,
          "linkTooltip": "Drill down",
          "linkUrl": "",
          "pattern": "Value #C",
          "thresholds": [],
          "type": "number",
          "unit": "percentunit"
        },
        {
          "alias": "CPU Limits",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": false,
          "linkTooltip": "Drill down",
          "linkUrl": "",
          "pattern": "Value #D",
          "thresholds": [],
          "type": "number",
          "unit": "short"
        },
        {
          "alias": "CPU Limits %",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": false,
          "linkTooltip": "Drill down",
          "linkUrl": "",
          "pattern": "Value #E",
          "thresholds": [],
          "type": "number",
          "unit": "percentunit"
        },
        {
          "alias": "Pod",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": true,
          "linkTooltip": "Drill down",
          "linkUrl": "/d/6581e46e4e5c7ba40a07646395ef7b55/kubernetes-compute-resources-pod?var-namespace=$namespace&var-pod=$__cell",
          "pattern": "pod",
          "thresholds": [],
          "type": "number",
          "unit": "short"
        },
        {
          "alias": "",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "pattern": "/.*/",
          "thresholds": [],
          "type": "string",
          "unit": "short"
        }
      ],
      "targets": [
        {
          "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
          "format": "table",
          "instant": true,
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A",
          "step": 10
        },
        {
          "expr": "sum(kube_pod_container_resource_requests{unit='core',namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
          "format": "table",
          "instant": true,
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "B",
          "step": 10
        },
        {
          "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod) / sum(kube_pod_container_resource_requests{unit='core',namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
          "format": "table",
          "instant": true,
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "C",
          "step": 10
        },
        {
          "expr": "sum(kube_pod_container_resource_limits{unit='core',namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
          "format": "table",
          "instant": true,
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "D",
          "step": 10
        },
        {
          "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod) / sum(kube_pod_container_resource_limits{unit='core',namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
          "format": "table",
          "instant": true,
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "E",
          "step": 10
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "CPU Quota",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "transform": "table",
      "type": "table-old",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": 0,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": false
        }
      ]
    },
    {
      "collapsed": false,
      "datasource": null,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 162
      },
      "id": 6,
      "panels": [],
      "repeat": "deployment",
      "scopedVars": {
        "deployment": {}
      },
      "title": "Memory Usage ($deployment)",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 7,
        "w": 24,
        "x": 0,
        "y": 163
      },
      "hiddenSeries": false,
      "id": 2,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null as zero",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace', pod=~'$deployment-[a-z0-9]+-[a-z0-9]+', container!=''}) by (pod)",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "legendLink": null,
          "refId": "A",
          "step": 10
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Memory Usage",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "bytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": 0,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": false
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "collapsed": false,
      "datasource": null,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 170
      },
      "id": 7,
      "panels": [],
      "repeat": "deployment",
      "scopedVars": {
        "deployment": {}
      },
      "title": "Memory Quota ($deployment)",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "columns": [],
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fontSize": "100%",
      "gridPos": {
        "h": 7,
        "w": 24,
        "x": 0,
        "y": 171
      },
      "id": 3,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null as zero",
      "pageSize": null,
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "scroll": true,
      "seriesOverrides": [],
      "showHeader": true,
      "sort": {
        "col": 1,
        "desc": true
      },
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "styles": [
        {
          "alias": "Time",
          "align": "auto",
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "pattern": "Time",
          "type": "hidden"
        },
        {
          "alias": "Memory Usage",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": false,
          "linkTooltip": "Drill down",
          "linkUrl": "",
          "pattern": "Value #A",
          "thresholds": [],
          "type": "number",
          "unit": "decbytes"
        },
        {
          "alias": "Memory Requests",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": false,
          "linkTooltip": "Drill down",
          "linkUrl": "",
          "pattern": "Value #B",
          "thresholds": [],
          "type": "number",
          "unit": "decbytes"
        },
        {
          "alias": "Memory Requests %",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": false,
          "linkTooltip": "Drill down",
          "linkUrl": "",
          "pattern": "Value #C",
          "thresholds": [],
          "type": "number",
          "unit": "percentunit"
        },
        {
          "alias": "Memory Limits",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": false,
          "linkTooltip": "Drill down",
          "linkUrl": "",
          "pattern": "Value #D",
          "thresholds": [],
          "type": "number",
          "unit": "decbytes"
        },
        {
          "alias": "Memory Limits %",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": false,
          "linkTooltip": "Drill down",
          "linkUrl": "",
          "pattern": "Value #E",
          "thresholds": [],
          "type": "number",
          "unit": "percentunit"
        },
        {
          "alias": "Pod",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "link": true,
          "linkTooltip": "Drill down",
          "linkUrl": "/d/6581e46e4e5c7ba40a07646395ef7b55/kubernetes-compute-resources-pod?var-namespace=$namespace&var-pod=$__cell",
          "pattern": "pod",
          "thresholds": [],
          "type": "number",
          "unit": "short"
        },
        {
          "alias": "",
          "align": "auto",
          "colorMode": null,
          "colors": [],
          "dateFormat": "YYYY-MM-DD HH:mm:ss",
          "decimals": 2,
          "pattern": "/.*/",
          "thresholds": [],
          "type": "string",
          "unit": "short"
        }
      ],
      "targets": [
        {
          "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+', container!=''}) by (pod)",
          "format": "table",
          "instant": true,
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A",
          "step": 10
        },
        {
          "expr": "sum(kube_pod_container_resource_requests{unit='byte',namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
          "format": "table",
          "instant": true,
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "B",
          "step": 10
        },
        {
          "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+', container!=''}) by (pod) / sum(kube_pod_container_resource_requests{unit='byte',namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
          "format": "table",
          "instant": true,
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "C",
          "step": 10
        },
        {
          "expr": "sum(kube_pod_container_resource_limits{unit='byte',namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
          "format": "table",
          "instant": true,
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "D",
          "step": 10
        },
        {
          "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+', container!=''}) by (pod) / sum(kube_pod_container_resource_limits{unit='byte',namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
          "format": "table",
          "instant": true,
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "E",
          "step": 10
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeShift": null,
      "title": "Memory Quota",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "transform": "table",
      "type": "table-old",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": 0,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": false
        }
      ]
    },
    {
      "collapsed": false,
      "datasource": null,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 178
      },
      "id": 15,
      "panels": [],
      "repeat": "deployment",
      "scopedVars": {
        "deployment": {}
      },
      "title": "Network Usage ($deployment)",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 6,
        "w": 24,
        "x": 0,
        "y": 179
      },
      "hiddenSeries": false,
      "id": 17,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(irate(container_network_receive_bytes_total{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}[5m])) by (pod)",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Receive Bandwidth",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "decimals": null,
          "format": "Bps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": false
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "$datasource",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 6,
        "w": 24,
        "x": 0,
        "y": 185
      },
      "hiddenSeries": false,
      "id": 18,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.5.15",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "scopedVars": {
        "deployment": {}
      },
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(irate(container_network_transmit_bytes_total{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}[5m])) by (pod)",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "{{`{{pod}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Transmit Bandwidth",
      "tooltip": {
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "Bps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": false
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    }
  ],
  "refresh": "10s",
  "schemaVersion": 27,
  "style": "dark",
  "tags": [
    "3scale",
    "redis",
    "twemproxy"
  ],
  "templating": {
    "list": [
      {
        "description": null,
        "error": null,
        "hide": 0,
        "includeAll": false,
        "label": null,
        "multi": false,
        "name": "datasource",
        "options": [],
        "query": "prometheus",
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "type": "datasource"
      },
      {
        "allValue": null,
        "current": {
          "selected": false,
          "text": "{{ .Namespace }}",
          "value": "{{ .Namespace }}"
        },
        "description": null,
        "error": null,
        "hide": 0,
        "includeAll": false,
        "label": "namespace",
        "multi": false,
        "name": "namespace",
        "options": [
          {
            "selected": true,
            "text": "{{ .Namespace }}",
            "value": "{{ .Namespace }}"
          }
        ],
        "query": "{{ .Namespace }}",
        "skipUrlSync": false,
        "type": "custom"
      },
      {
        "allValue": null,
        "datasource": "$datasource",
        "definition": "label_values(twemproxy_exporter_up{namespace='$namespace',}, pod)",
        "description": null,
        "error": null,
        "hide": 0,
        "includeAll": true,
        "label": "deployment",
        "multi": false,
        "name": "deployment",
        "options": [],
        "query": {
          "query": "label_values(twemproxy_exporter_up{namespace='$namespace',}, pod)",
          "refId": "StandardVariableQuery"
        },
        "refresh": 1,
        "regex": "/(.*)-[a-z0-9]+-[a-z0-9]+/",
        "skipUrlSync": false,
        "sort": 1,
        "tagValuesQuery": "",
        "tags": [],
        "tagsQuery": "",
        "type": "query",
        "useTags": false
      },
      {
        "allValue": null,
        "datasource": "${datasource}",
        "definition": "label_values(twemproxy_exporter_client_err_total{namespace='$namespace'},pool)",
        "description": null,
        "error": null,
        "hide": 0,
        "includeAll": false,
        "label": null,
        "multi": false,
        "name": "pool",
        "options": [],
        "query": {
          "query": "label_values(twemproxy_exporter_client_err_total{namespace='$namespace'},pool)",
          "refId": "StandardVariableQuery"
        },
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 0,
        "tagValuesQuery": "",
        "tags": [],
        "tagsQuery": "",
        "type": "query",
        "useTags": false
      }
    ]
  },
  "time": {
    "from": "now-1h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": [
      "10s",
      "30s",
      "1m",
      "5m",
      "15m",
      "30m",
      "1h",
      "2h",
      "1d"
    ],
    "time_options": [
      "5m",
      "15m",
      "1h",
      "6h",
      "12h",
      "24h",
      "2d",
      "7d",
      "30d"
    ]
  },
  "timezone": "",
  "title": "3scale Redis Twemproxy {{ .Name }}"
}
