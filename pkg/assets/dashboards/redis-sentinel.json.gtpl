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
    "graphTooltip": 1,
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
        "repeat": "sentinel",
        "title": "Sentinel [[sentinel]]",
        "type": "row"
      },
      {
        "aliasColors": {},
        "bars": true,
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
          "w": 10,
          "x": 0,
          "y": 1
        },
        "hiddenSeries": false,
        "id": 42,
        "legend": {
          "alignAsTable": false,
          "avg": false,
          "current": false,
          "max": false,
          "min": false,
          "rightSide": false,
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "exemplar": true,
            "expr": "rate(saas_redis_sentinel_sdown_count{namespace='$namespace',sentinel=~\"[[sentinel]].*\"}[1m]) > 0",
            "format": "time_series",
            "interval": "",
            "intervalFactor": 1,
            "legendFormat": "{{`{{shard}}`}}-{{`{{redis_server}}`}}",
            "refId": "A"
          },
          {
            "exemplar": true,
            "expr": "rate(saas_redis_sentinel_sdown_sentinel_count{namespace='$namespace',sentinel=~\"[[sentinel]].*\"}[1m]) > 0",
            "format": "time_series",
            "interval": "",
            "intervalFactor": 1,
            "legendFormat": "sentinel-{{`{{redis_server}}`}}",
            "refId": "B"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Subjective Down",
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
            "format": "none",
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
        "bars": true,
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
          "w": 10,
          "x": 10,
          "y": 1
        },
        "hiddenSeries": false,
        "id": 41,
        "legend": {
          "alignAsTable": false,
          "avg": false,
          "current": false,
          "max": false,
          "min": false,
          "rightSide": false,
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "exemplar": true,
            "expr": "rate(saas_redis_sentinel_switch_master_count{namespace='$namespace',sentinel=~\"[[sentinel]].*\"}[1m]) > 0",
            "format": "time_series",
            "interval": "",
            "intervalFactor": 1,
            "legendFormat": "{{`{{shard}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Failover",
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
            "format": "none",
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
        "bars": true,
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
          "w": 4,
          "x": 20,
          "y": 1
        },
        "hiddenSeries": false,
        "id": 43,
        "legend": {
          "alignAsTable": false,
          "avg": false,
          "current": false,
          "max": false,
          "min": false,
          "rightSide": false,
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "exemplar": true,
            "expr": "rate(saas_redis_sentinel_failover_abort_no_good_slave_count{namespace='$namespace',sentinel=~\"[[sentinel]].*\"}[1m]) > 0",
            "format": "time_series",
            "interval": "",
            "intervalFactor": 1,
            "legendFormat": "{{`{{shard}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Failover Abort No Good Slave",
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
            "format": "none",
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
          "y": 9
        },
        "hiddenSeries": false,
        "id": 22,
        "legend": {
          "alignAsTable": false,
          "avg": false,
          "current": false,
          "max": false,
          "min": false,
          "rightSide": true,
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "exemplar": true,
            "expr": "saas_redis_sentinel_last_ok_ping_reply{namespace='$namespace',sentinel=~\"[[sentinel]].*\"}/1000",
            "format": "time_series",
            "interval": "",
            "intervalFactor": 1,
            "legendFormat": "{{`{{shard}}`}}-{{`{{redis_server}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [
          {
            "colorMode": "critical",
            "fill": true,
            "line": true,
            "op": "gt",
            "value": null,
            "yaxis": "left"
          }
        ],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Last Ping OK reply",
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
            "format": "s",
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
        "id": 26,
        "legend": {
          "alignAsTable": false,
          "avg": false,
          "current": false,
          "max": false,
          "min": false,
          "rightSide": true,
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "exemplar": true,
            "expr": "saas_redis_sentinel_master_link_down_time{namespace='$namespace',sentinel=~\"[[sentinel]].*\"} ",
            "format": "time_series",
            "interval": "",
            "intervalFactor": 1,
            "legendFormat": "{{`{{shard}}`}}-{{`{{redis_server}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Master Link Down",
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
            "min": "0",
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
        "id": 28,
        "legend": {
          "avg": false,
          "current": false,
          "max": false,
          "min": false,
          "rightSide": true,
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "exemplar": true,
            "expr": "saas_redis_sentinel_link_pending_commands{namespace='$namespace',sentinel=~\"[[sentinel]].*\"}",
            "format": "time_series",
            "interval": "",
            "intervalFactor": 1,
            "legendFormat": "{{`{{shard}}`}}-{{`{{redis_server}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Link Pending Commands",
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
            "min": "0",
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
          "y": 17
        },
        "hiddenSeries": false,
        "id": 40,
        "legend": {
          "alignAsTable": false,
          "avg": false,
          "current": false,
          "max": false,
          "min": false,
          "rightSide": true,
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "exemplar": true,
            "expr": "saas_redis_sentinel_slave_repl_offset{namespace='$namespace',sentinel=~\"[[sentinel]].*\"}",
            "format": "time_series",
            "interval": "",
            "intervalFactor": 1,
            "legendFormat": "{{`{{shard}}`}}-{{`{{redis_server}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Slave Replication Offset",
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
            "format": "none",
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
          "w": 8,
          "x": 0,
          "y": 25
        },
        "hiddenSeries": false,
        "id": 38,
        "legend": {
          "alignAsTable": false,
          "avg": false,
          "current": false,
          "max": false,
          "min": false,
          "rightSide": false,
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "exemplar": true,
            "expr": "saas_redis_sentinel_num_slaves{namespace='$namespace',sentinel=~\"[[sentinel]].*\"} ",
            "format": "time_series",
            "interval": "",
            "intervalFactor": 1,
            "legendFormat": "{{`{{shard}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Num Slaves",
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
            "min": "0",
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
        "description": "",
        "fieldConfig": {
          "defaults": {},
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 8,
          "x": 8,
          "y": 25
        },
        "hiddenSeries": false,
        "id": 24,
        "legend": {
          "avg": false,
          "current": false,
          "max": false,
          "min": false,
          "rightSide": false,
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "exemplar": true,
            "expr": "saas_redis_sentinel_num_other_sentinels{namespace='$namespace',sentinel=~\"[[sentinel]].*\"}",
            "format": "time_series",
            "interval": "",
            "intervalFactor": 1,
            "legendFormat": "{{`{{shard}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Num Other Sentinels",
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
        "description": "",
        "fieldConfig": {
          "defaults": {},
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 8,
          "x": 16,
          "y": 25
        },
        "hiddenSeries": false,
        "id": 39,
        "legend": {
          "alignAsTable": false,
          "avg": false,
          "current": false,
          "max": false,
          "min": false,
          "rightSide": false,
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "exemplar": true,
            "expr": "saas_redis_sentinel_role_reported_time{namespace='$namespace',sentinel=~\"[[sentinel]].*\"}",
            "format": "time_series",
            "interval": "",
            "intervalFactor": 1,
            "legendFormat": "{{`{{shard}}`}}-{{`{{redis_server}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Role Reported time",
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
            "format": "ms",
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
          "y": 33
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
          "y": 34
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
        "sparkline": {
          "fillColor": "rgba(31, 118, 189, 0.18)",
          "full": false,
          "lineColor": "rgb(31, 120, 193)",
          "show": false
        },
        "tableColumn": "",
        "targets": [
          {
            "expr": "kube_statefulset_status_replicas_ready{namespace='$namespace',statefulset=~'redis-sentinel.*'}",
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
          "y": 34
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
        "sparkline": {
          "fillColor": "rgba(31, 118, 189, 0.18)",
          "full": false,
          "lineColor": "rgb(31, 120, 193)",
          "show": false
        },
        "tableColumn": "",
        "targets": [
          {
            "expr": "kube_statefulset_status_replicas{namespace='$namespace',statefulset=~'redis-sentinel.*'} - kube_statefulset_status_replicas_ready{namespace='$namespace',statefulset=~'redis-sentinel.*'}",
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
          "y": 34
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
        "sparkline": {
          "fillColor": "rgba(31, 118, 189, 0.18)",
          "full": false,
          "lineColor": "rgb(31, 120, 193)",
          "show": false
        },
        "tableColumn": "",
        "targets": [
          {
            "expr": "count(count(container_memory_working_set_bytes{namespace='$namespace',pod=~'redis-sentinel.*'}) by (node))",
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
          "y": 34
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
        "sparkline": {
          "fillColor": "rgba(31, 118, 189, 0.18)",
          "full": false,
          "lineColor": "rgb(31, 120, 193)",
          "show": false
        },
        "tableColumn": "",
        "targets": [
          {
            "expr": "max(sum(delta(kube_pod_container_status_restarts_total{namespace='$namespace',pod=~'redis-sentinel.*'}[5m])) by (pod))",
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
          "y": 37
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "kube_statefulset_status_replicas{namespace='$namespace',statefulset=~'redis-sentinel.*'}",
            "format": "time_series",
            "intervalFactor": 2,
            "legendFormat": "total-pods",
            "legendLink": null,
            "refId": "A",
            "step": 10
          },
          {
            "expr": "kube_statefulset_status_replicas_ready{namespace='$namespace',statefulset=~'redis-sentinel.*'}",
            "format": "time_series",
            "intervalFactor": 1,
            "legendFormat": "avail-pods",
            "refId": "B"
          },
          {
            "expr": "kube_statefulset_status_replicas {namespace='$namespace',statefulset=~'redis-sentinel.*'} - kube_statefulset_status_replicas_ready {namespace='$namespace',statefulset=~'redis-sentinel.*'}",
            "format": "time_series",
            "intervalFactor": 1,
            "legendFormat": "unavail-pods",
            "refId": "C"
          },
          {
            "expr": "count(count(container_memory_working_set_bytes{namespace='$namespace',pod=~'redis-sentinel.*'}) by (node))",
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
          "y": 44
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(delta(kube_pod_container_status_restarts_total{namespace='$namespace',pod=~'redis-sentinel.*'}[5m])) by (pod)",
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
          "y": 50
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
          "y": 51
        },
        "hiddenSeries": false,
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
        "options": {
          "alertThreshold": true
        },
        "percentage": false,
        "pluginVersion": "7.5.15",
        "pointradius": 5,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod)",
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
          "y": 58
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
          "y": 59
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
            "linkUrl": "/d/6581e46e4e5c7ba40a07646395ef7b55/3scale-compute-resources-pod?var-namespace=$namespace&var-pod=$__cell",
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
            "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod)",
            "format": "table",
            "instant": true,
            "intervalFactor": 2,
            "legendFormat": "",
            "refId": "A",
            "step": 10
          },
          {
            "expr": "sum(kube_pod_container_resource_requests{unit='core',namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod)",
            "format": "table",
            "instant": true,
            "intervalFactor": 2,
            "legendFormat": "",
            "refId": "B",
            "step": 10
          },
          {
            "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod) / sum(kube_pod_container_resource_requests{unit='core',namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod)",
            "format": "table",
            "instant": true,
            "intervalFactor": 2,
            "legendFormat": "",
            "refId": "C",
            "step": 10
          },
          {
            "expr": "sum(kube_pod_container_resource_limits{unit='core',namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod)",
            "format": "table",
            "instant": true,
            "intervalFactor": 2,
            "legendFormat": "",
            "refId": "D",
            "step": 10
          },
          {
            "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod) / sum(kube_pod_container_resource_limits{unit='core',namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod)",
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
          "y": 66
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
          "y": 67
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace', pod=~'redis-sentinel.*', container!=''}) by (pod)",
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
          "y": 74
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
          "y": 75
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
            "linkUrl": "/d/6581e46e4e5c7ba40a07646395ef7b55/3scale-compute-resources-pod?var-namespace=$namespace&var-pod=$__cell",
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
            "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace', pod=~'redis-sentinel.*', container!=''}) by (pod)",
            "format": "table",
            "instant": true,
            "intervalFactor": 2,
            "legendFormat": "",
            "refId": "A",
            "step": 10
          },
          {
            "expr": "sum(kube_pod_container_resource_requests{unit='byte',namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod)",
            "format": "table",
            "instant": true,
            "intervalFactor": 2,
            "legendFormat": "",
            "refId": "B",
            "step": 10
          },
          {
            "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace', pod=~'redis-sentinel.*', container!=''}) by (pod) / sum(kube_pod_container_resource_requests{unit='byte',namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod)",
            "format": "table",
            "instant": true,
            "intervalFactor": 2,
            "legendFormat": "",
            "refId": "C",
            "step": 10
          },
          {
            "expr": "sum(kube_pod_container_resource_limits{unit='byte',namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod)",
            "format": "table",
            "instant": true,
            "intervalFactor": 2,
            "legendFormat": "",
            "refId": "D",
            "step": 10
          },
          {
            "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace', pod=~'redis-sentinel.*', container!=''}) by (pod) / sum(kube_pod_container_resource_limits{unit='byte',namespace=~'$namespace', pod=~'redis-sentinel.*'}) by (pod)",
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
          "y": 82
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
          "y": 83
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(irate(container_network_receive_bytes_total{namespace=~'$namespace', pod=~'redis-sentinel.*'}[5m])) by (pod)",
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
          "y": 89
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
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(irate(container_network_transmit_bytes_total{namespace=~'$namespace', pod=~'redis-sentinel.*'}[5m])) by (pod)",
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
    "schemaVersion": 27,
    "style": "dark",
    "tags": [
        "3scale",
        "redis",
        "sentinel"
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
          "current": {
            "selected": false,
            "text": "All",
            "value": "$__all"
          },
          "datasource": "${datasource}",
          "definition": "label_values(saas_redis_sentinel_last_ok_ping_reply{namespace='$namespace'},sentinel)",
          "description": null,
          "error": null,
          "hide": 0,
          "includeAll": true,
          "label": null,
          "multi": true,
          "name": "sentinel",
          "options": [],
          "query": {
            "query": "label_values(saas_redis_sentinel_last_ok_ping_reply{namespace='$namespace'},sentinel)",
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
    "title": "3scale Redis Sentinel"
  }