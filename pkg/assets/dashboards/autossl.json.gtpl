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
            "gridPos": {
                "h": 1,
                "w": 24,
                "x": 0,
                "y": 0
            },
            "id": 20,
            "panels": [],
            "title": "Application",
            "type": "row"
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": "$datasource",
            "fill": 1,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 1
            },
            "id": 22,
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
            "options": {},
            "percentage": false,
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [
                {
                    "alias": "error",
                    "color": "#C4162A"
                }
            ],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "expr": "sum(increase(autossl_domains_allowed{namespace='$namespace', pod=~'autossl.*'}[24h])) by (status)",
                    "format": "time_series",
                    "intervalFactor": 1,
                    "legendFormat": "{{`{{state}}`}}",
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeRegions": [],
            "timeShift": null,
            "title": "Domain validations (last 24h)",
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
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": "$datasource",
            "fill": 1,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 1
            },
            "id": 24,
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
            "options": {},
            "percentage": false,
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [
                {
                    "alias": "error",
                    "color": "#C4162A"
                },
                {
                    "alias": "ok",
                    "color": "#37872D"
                }
            ],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "expr": "sum(increase(autossl_letsencrypt_requests{namespace='$namespace', pod=~'autossl.*'}[24h])) by (status)",
                    "format": "time_series",
                    "intervalFactor": 1,
                    "legendFormat": "{{`{{state}}`}}",
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeRegions": [],
            "timeShift": null,
            "title": "Let's encrypt issued certificates (last 24h)",
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
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": "$datasource",
            "fill": 1,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 0,
                "y": 9
            },
            "id": 26,
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
            "options": {},
            "percentage": false,
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "expr": "sum(nginx_http_connections{namespace='$namespace', pod=~'autossl.*'}) by (state)",
                    "format": "time_series",
                    "intervalFactor": 1,
                    "legendFormat": "{{`{{state}}`}}",
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeRegions": [],
            "timeShift": null,
            "title": "Autossl Connections",
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
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": "$datasource",
            "fill": 1,
            "gridPos": {
                "h": 8,
                "w": 12,
                "x": 12,
                "y": 9
            },
            "id": 28,
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
            "options": {},
            "percentage": false,
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [
                {
                    "alias": "errors",
                    "color": "#C4162A"
                }
            ],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "expr": "sum(nginx_metric_errors_total{namespace='$namespace', pod=~'autossl.*'})",
                    "format": "time_series",
                    "intervalFactor": 1,
                    "legendFormat": "errors",
                    "refId": "A"
                }
            ],
            "thresholds": [],
            "timeFrom": null,
            "timeRegions": [],
            "timeShift": null,
            "title": "Prometheus Metrics Errors",
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
            "gridPos": {
                "h": 1,
                "w": 24,
                "x": 0,
                "y": 17
            },
            "id": 13,
            "panels": [],
            "title": "Pods",
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
                "y": 18
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
            "options": {},
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
            "sparkline": {
                "fillColor": "rgba(31, 118, 189, 0.18)",
                "full": false,
                "lineColor": "rgb(31, 120, 193)",
                "show": false
            },
            "tableColumn": "",
            "targets": [
                {
                    "expr": "kube_deployment_status_replicas_available{namespace='$namespace',deployment=~'autossl.*'}",
                    "format": "time_series",
                    "intervalFactor": 1,
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
                "#d44a3a"
            ],
            "datasource": "$datasource",
            "decimals": 0,
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
                "y": 18
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
            "options": {},
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
            "sparkline": {
                "fillColor": "rgba(31, 118, 189, 0.18)",
                "full": false,
                "lineColor": "rgb(31, 120, 193)",
                "show": false
            },
            "tableColumn": "",
            "targets": [
                {
                    "expr": "kube_deployment_status_replicas_unavailable{namespace='$namespace',deployment=~'autossl.*'}",
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
                "y": 18
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
            "options": {},
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
            "sparkline": {
                "fillColor": "rgba(31, 118, 189, 0.18)",
                "full": false,
                "lineColor": "rgb(31, 120, 193)",
                "show": false
            },
            "tableColumn": "",
            "targets": [
                {
                    "expr": "count(count(container_memory_working_set_bytes{namespace='$namespace',pod=~'autossl.*'}) by (node))",
                    "format": "time_series",
                    "intervalFactor": 1,
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
                "y": 18
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
            "options": {},
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
            "sparkline": {
                "fillColor": "rgba(31, 118, 189, 0.18)",
                "full": false,
                "lineColor": "rgb(31, 120, 193)",
                "show": false
            },
            "tableColumn": "",
            "targets": [
                {
                    "expr": "max(sum(delta(kube_pod_container_status_restarts_total{namespace='$namespace',pod=~'autossl.*'}[5m])) by (pod))",
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
            "fill": 1,
            "gridPos": {
                "h": 7,
                "w": 24,
                "x": 0,
                "y": 21
            },
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
            "options": {},
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "expr": "kube_deployment_status_replicas{namespace='$namespace',deployment=~'autossl.*'}",
                    "format": "time_series",
                    "intervalFactor": 2,
                    "legendFormat": "total-pods",
                    "legendLink": null,
                    "refId": "A",
                    "step": 10
                },
                {
                    "expr": "kube_deployment_status_replicas_available{namespace='$namespace',deployment=~'autossl.*'}",
                    "format": "time_series",
                    "intervalFactor": 1,
                    "legendFormat": "avail-pods",
                    "refId": "B"
                },
                {
                    "expr": "kube_deployment_status_replicas_unavailable{namespace='$namespace',deployment=~'autossl.*'}",
                    "format": "time_series",
                    "intervalFactor": 1,
                    "legendFormat": "unavail-pods",
                    "refId": "C"
                },
                {
                    "expr": "count(count(container_memory_working_set_bytes{namespace='$namespace',pod=~'autossl.*'}) by (node))",
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
            "fill": 1,
            "gridPos": {
                "h": 6,
                "w": 24,
                "x": 0,
                "y": 28
            },
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
            "options": {},
            "percentage": false,
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "expr": "sum(delta(kube_pod_container_status_restarts_total{namespace='$namespace',pod=~'autossl.*'}[5m])) by (pod)",
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
            "gridPos": {
                "h": 1,
                "w": 24,
                "x": 0,
                "y": 34
            },
            "id": 4,
            "panels": [],
            "repeat": null,
            "title": "CPU Usage",
            "type": "row"
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": "$datasource",
            "fill": 1,
            "gridPos": {
                "h": 7,
                "w": 24,
                "x": 0,
                "y": 35
            },
            "id": 0,
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
            "options": {},
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=~'$namespace', pod=~'autossl.*'}) by (pod)",
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
            "gridPos": {
                "h": 1,
                "w": 24,
                "x": 0,
                "y": 42
            },
            "id": 5,
            "panels": [],
            "repeat": null,
            "title": "CPU Quota",
            "type": "row"
        },
        {
            "aliasColors": {},
            "bars": false,
            "columns": [],
            "dashLength": 10,
            "dashes": false,
            "datasource": "$datasource",
            "fill": 1,
            "fontSize": "100%",
            "gridPos": {
                "h": 7,
                "w": 24,
                "x": 0,
                "y": 43
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
            "options": {},
            "pageSize": null,
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
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
                    "dateFormat": "YYYY-MM-DD HH:mm:ss",
                    "pattern": "Time",
                    "type": "hidden"
                },
                {
                    "alias": "CPU Usage",
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
                    "colorMode": null,
                    "colors": [],
                    "dateFormat": "YYYY-MM-DD HH:mm:ss",
                    "decimals": 2,
                    "link": true,
                    "linkTooltip": "Drill down",
                    "linkUrl": "/d/6581e46e4e5c7ba40a07646395ef7b55/3scale-compute-resources-pod?var-namespace=$namespace&var-pod=$__cell",
                    "pattern": "pod",
                    "thresholds": [],
                    "type": "number",
                    "unit": "short"
                },
                {
                    "alias": "",
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
                    "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=~'$namespace', pod=~'autossl.*'}) by (pod)",
                    "format": "table",
                    "instant": true,
                    "intervalFactor": 2,
                    "legendFormat": "",
                    "refId": "A",
                    "step": 10
                },
                {
                    "expr": "sum(kube_pod_container_resource_requests{unit='core',namespace=~'$namespace', pod=~'autossl.*'}) by (pod)",
                    "format": "table",
                    "instant": true,
                    "intervalFactor": 2,
                    "legendFormat": "",
                    "refId": "B",
                    "step": 10
                },
                {
                    "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=~'$namespace', pod=~'autossl.*'}) by (pod) / sum(kube_pod_container_resource_requests{unit='core',namespace=~'$namespace', pod=~'autossl.*'}) by (pod)",
                    "format": "table",
                    "instant": true,
                    "intervalFactor": 2,
                    "legendFormat": "",
                    "refId": "C",
                    "step": 10
                },
                {
                    "expr": "sum(kube_pod_container_resource_limits{unit='core',namespace=~'$namespace', pod=~'autossl.*'}) by (pod)",
                    "format": "table",
                    "instant": true,
                    "intervalFactor": 2,
                    "legendFormat": "",
                    "refId": "D",
                    "step": 10
                },
                {
                    "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace=~'$namespace', pod=~'autossl.*'}) by (pod) / sum(kube_pod_container_resource_limits{unit='core',namespace=~'$namespace', pod=~'autossl.*'}) by (pod)",
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
            "type": "table",
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
            "gridPos": {
                "h": 1,
                "w": 24,
                "x": 0,
                "y": 50
            },
            "id": 6,
            "panels": [],
            "repeat": null,
            "title": "Memory Usage",
            "type": "row"
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": "$datasource",
            "fill": 1,
            "gridPos": {
                "h": 7,
                "w": 24,
                "x": 0,
                "y": 51
            },
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
            "options": {},
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace', pod=~'autossl.*', container!=''}) by (pod)",
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
            "gridPos": {
                "h": 1,
                "w": 24,
                "x": 0,
                "y": 58
            },
            "id": 7,
            "panels": [],
            "repeat": null,
            "title": "Memory Quota",
            "type": "row"
        },
        {
            "aliasColors": {},
            "bars": false,
            "columns": [],
            "dashLength": 10,
            "dashes": false,
            "datasource": "$datasource",
            "fill": 1,
            "fontSize": "100%",
            "gridPos": {
                "h": 7,
                "w": 24,
                "x": 0,
                "y": 59
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
            "options": {},
            "pageSize": null,
            "percentage": false,
            "pointradius": 5,
            "points": false,
            "renderer": "flot",
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
                    "dateFormat": "YYYY-MM-DD HH:mm:ss",
                    "pattern": "Time",
                    "type": "hidden"
                },
                {
                    "alias": "Memory Usage",
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
                    "colorMode": null,
                    "colors": [],
                    "dateFormat": "YYYY-MM-DD HH:mm:ss",
                    "decimals": 2,
                    "link": true,
                    "linkTooltip": "Drill down",
                    "linkUrl": "/d/6581e46e4e5c7ba40a07646395ef7b55/3scale-compute-resources-pod?var-namespace=$namespace&var-pod=$__cell",
                    "pattern": "pod",
                    "thresholds": [],
                    "type": "number",
                    "unit": "short"
                },
                {
                    "alias": "",
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
                    "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace', pod=~'autossl.*', container!=''}) by (pod)",
                    "format": "table",
                    "instant": true,
                    "intervalFactor": 2,
                    "legendFormat": "",
                    "refId": "A",
                    "step": 10
                },
                {
                    "expr": "sum(kube_pod_container_resource_requests{unit='byte',namespace=~'$namespace', pod=~'autossl.*'}) by (pod)",
                    "format": "table",
                    "instant": true,
                    "intervalFactor": 2,
                    "legendFormat": "",
                    "refId": "B",
                    "step": 10
                },
                {
                    "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace', pod=~'autossl.*', container!=''}) by (pod) / sum(kube_pod_container_resource_requests{unit='byte',namespace=~'$namespace', pod=~'autossl.*'}) by (pod)",
                    "format": "table",
                    "instant": true,
                    "intervalFactor": 2,
                    "legendFormat": "",
                    "refId": "C",
                    "step": 10
                },
                {
                    "expr": "sum(kube_pod_container_resource_limits{unit='byte',namespace=~'$namespace', pod=~'autossl.*'}) by (pod)",
                    "format": "table",
                    "instant": true,
                    "intervalFactor": 2,
                    "legendFormat": "",
                    "refId": "D",
                    "step": 10
                },
                {
                    "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace', pod=~'autossl.*', container!=''}) by (pod) / sum(kube_pod_container_resource_limits{unit='byte',namespace=~'$namespace', pod=~'autossl.*'}) by (pod)",
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
            "type": "table",
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
            "gridPos": {
                "h": 1,
                "w": 24,
                "x": 0,
                "y": 66
            },
            "id": 15,
            "panels": [],
            "title": "Network Usage",
            "type": "row"
        },
        {
            "aliasColors": {},
            "bars": false,
            "dashLength": 10,
            "dashes": false,
            "datasource": "$datasource",
            "fill": 1,
            "gridPos": {
                "h": 6,
                "w": 24,
                "x": 0,
                "y": 67
            },
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
            "options": {},
            "percentage": false,
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "expr": "sum(irate(container_network_receive_bytes_total{namespace=~'$namespace', pod=~'autossl.*'}[5m])) by (pod)",
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
            "fill": 1,
            "gridPos": {
                "h": 6,
                "w": 24,
                "x": 0,
                "y": 73
            },
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
            "options": {},
            "percentage": false,
            "pointradius": 2,
            "points": false,
            "renderer": "flot",
            "seriesOverrides": [],
            "spaceLength": 10,
            "stack": false,
            "steppedLine": false,
            "targets": [
                {
                    "expr": "sum(irate(container_network_transmit_bytes_total{namespace=~'$namespace', pod=~'autossl.*'}[5m])) by (pod)",
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
        }
    ],
    "refresh": "10s",
    "schemaVersion": 18,
    "style": "dark",
    "tags": [
        "3scale",
        "autossl",
        "system"
    ],
    "templating": {
        "list": [
            {
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
                    "tags": [],
                    "text": "{{ .Namespace }}",
                    "value": "{{ .Namespace }}"
                },
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
            }
        ]
    },
    "time": {
        "from": "now-1h",
        "to": "now"
    },
    "timepicker": {
        "refresh_intervals": [
            "5s",
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
    "title": "3scale System AutoSSL"
}