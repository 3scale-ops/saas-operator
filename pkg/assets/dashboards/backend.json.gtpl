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
        "title": "Listeners",
        "type": "row"
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 1
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_codes{namespace=\"$namespace\"}[1m])) by (request_type)",
            "format": "time_series",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{request_type}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Listener requests per second (by request type)",
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 1
        },
        "hiddenSeries": false,
        "id": 41,
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_codes{namespace=\"$namespace\"}[1m])) by (resp_code)",
            "format": "time_series",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{resp_code}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Listener requests per second (by response code)",
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
          "defaults": {
            "custom": {}
          },
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
        "id": 43,
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"authorize\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Authorize requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 9
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 45,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"authorize\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Authorize requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
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
        "id": 47,
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"authrep\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Authrep requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#5794F2",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 17
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 59,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"authrep\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Authrep requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
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
        "id": 50,
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"report\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Report requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#F2495C",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 25
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 58,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"report\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Report requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
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
        "id": 60,
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"authorize_oauth\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Authorize (oauth) requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FF9830",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 33
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 61,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"authorize_oauth\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Authorize (oauth) requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
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
        "id": 62,
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"authrep_oauth\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Authrep (oauth) requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#B877D9",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 41
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 63,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"authrep_oauth\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Authrep (oauth) requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "collapsed": false,
        "datasource": null,
        "gridPos": {
          "h": 1,
          "w": 24,
          "x": 0,
          "y": 49
        },
        "id": 66,
        "panels": [],
        "title": "Listeners Internal API",
        "type": "row"
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 50
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_codes{namespace=\"$namespace\"}[1m])) by (request_type)",
            "format": "time_series",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{request_type}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Listener Internal API requests per second (by request type)",
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 50
        },
        "hiddenSeries": false,
        "id": 68,
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_codes{namespace=\"$namespace\"}[1m])) by (resp_code)",
            "format": "time_series",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{resp_code}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Listener Internal API requests per second (by response code)",
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 58
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"alerts\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Alerts requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 58
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 70,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"alerts\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Alerts requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 66
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"application_keys\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Application Keys requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 66
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 72,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"application_keys\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Application Keys requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 74
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"application_referrer_filters\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Referrer Filters requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 74
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 74,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"application_referrer_filters\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Referrer Filters requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 82
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"applications\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Applications requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 82
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 76,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"applications\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Applications requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 90
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"errors\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Integration Errors API requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 90
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 78,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"errors\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Errors requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 98
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"events\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Events requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 98
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 80,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"events\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Events requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 106
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"metrics\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Metrics requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 106
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 82,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"metrics\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Metrics requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 114
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"service_tokens\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Service Tokens requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 114
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 84,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"service_tokens\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Service Tokens requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 122
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"services\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Services requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 122
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 86,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"services\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Services requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 130
        },
        "hiddenSeries": false,
        "id": 87,
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"stats\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Stats requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 130
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 88,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"stats\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Stats requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 138
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"usage_limits\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Usage Limits requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 138
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 90,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"usage_limits\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Usage Limits  requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 146
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"utilization\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Utilization requests per second (by response time bucket in seconds)",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#FADE2A",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "description": "",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 146
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 92,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_listener_internal_api_response_times_seconds_bucket{namespace=\"$namespace\",request_type=\"utilization\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Utilization  requests per second heatmap (by response time bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
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
        "id": 52,
        "panels": [],
        "title": "Workers",
        "type": "row"
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 24,
          "x": 0,
          "y": 155
        },
        "hiddenSeries": false,
        "id": 49,
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_worker_job_count{namespace=\"$namespace\"} [1m])) by (type)",
            "format": "time_series",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{type}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Jobs per second (by job type)",
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
            "format": "ops",
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 163
        },
        "hiddenSeries": false,
        "id": 54,
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_worker_job_runtime_seconds_bucket{namespace=\"$namespace\",type=\"ReportJob\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Report jobs per second (by runtime bucket in seconds)",
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
            "format": "ops",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#5794F2",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 163
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 55,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_worker_job_runtime_seconds_bucket{namespace=\"$namespace\",type=\"ReportJob\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Report jobs per second heatmap (by runtime bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "aliasColors": {},
        "bars": false,
        "dashLength": 10,
        "dashes": false,
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 0,
          "y": 171
        },
        "hiddenSeries": false,
        "id": 56,
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_worker_job_runtime_seconds_bucket{namespace=\"$namespace\",type=\"NotifyJob\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "thresholds": [],
        "timeFrom": null,
        "timeRegions": [],
        "timeShift": null,
        "title": "Notify jobs per second (by runtime bucket in seconds)",
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
            "format": "ops",
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
        "cards": {
          "cardPadding": null,
          "cardRound": null
        },
        "color": {
          "cardColor": "#b4ff00",
          "colorScale": "sqrt",
          "colorScheme": "interpolateOranges",
          "exponent": 0.5,
          "mode": "opacity"
        },
        "dataFormat": "tsbuckets",
        "datasource": "$datasource",
        "fieldConfig": {
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "gridPos": {
          "h": 8,
          "w": 12,
          "x": 12,
          "y": 171
        },
        "heatmap": {},
        "hideZeroBuckets": false,
        "highlightCards": true,
        "id": 57,
        "legend": {
          "show": false
        },
        "links": [],
        "reverseYBuckets": false,
        "targets": [
          {
            "expr": "sum(rate(apisonator_worker_job_runtime_seconds_bucket{namespace=\"$namespace\",type=\"NotifyJob\"}[1m])) by (le)",
            "format": "heatmap",
            "interval": "1m",
            "intervalFactor": 10,
            "legendFormat": "{{`{{le}}`}}",
            "refId": "A"
          }
        ],
        "timeFrom": null,
        "timeShift": null,
        "title": "Notify jobs per second heatmap (by runtime bucket in seconds)",
        "tooltip": {
          "show": true,
          "showHistogram": false
        },
        "type": "heatmap",
        "xAxis": {
          "show": true
        },
        "xBucketNumber": null,
        "xBucketSize": null,
        "yAxis": {
          "decimals": 0,
          "format": "s",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true,
          "splitFactor": null
        },
        "yBucketBound": "auto",
        "yBucketNumber": null,
        "yBucketSize": null
      },
      {
        "collapsed": false,
        "datasource": null,
        "gridPos": {
          "h": 1,
          "w": 24,
          "x": 0,
          "y": 179
        },
        "id": 13,
        "panels": [],
        "repeat": "deployment",
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
          "defaults": {
            "custom": {}
          },
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
          "y": 180
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
          "defaults": {
            "custom": {}
          },
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
          "y": 180
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
          "defaults": {
            "custom": {}
          },
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
          "y": 180
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
          "defaults": {
            "custom": {}
          },
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
          "y": 180
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 7,
          "w": 24,
          "x": 0,
          "y": 183
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 5,
        "points": false,
        "renderer": "flot",
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 6,
          "w": 24,
          "x": 0,
          "y": 190
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
        "repeatedByRow": false,
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
          "y": 196
        },
        "id": 4,
        "panels": [],
        "repeat": "deployment",
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 7,
          "w": 24,
          "x": 0,
          "y": 197
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 5,
        "points": false,
        "renderer": "flot",
        "seriesOverrides": [],
        "spaceLength": 10,
        "stack": false,
        "steppedLine": false,
        "targets": [
          {
            "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
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
          "y": 204
        },
        "id": 5,
        "panels": [],
        "repeat": "deployment",
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fontSize": "100%",
        "gridPos": {
          "h": 7,
          "w": 24,
          "x": 0,
          "y": 205
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
            "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
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
            "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod) / sum(kube_pod_container_resource_requests{unit='core',namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
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
            "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod) / sum(kube_pod_container_resource_limits{unit='core',namespace=~'$namespace',pod=~'$deployment-[a-z0-9]+-[a-z0-9]+'}) by (pod)",
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
          "y": 212
        },
        "id": 6,
        "panels": [],
        "repeat": "deployment",
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fillGradient": 0,
        "gridPos": {
          "h": 7,
          "w": 24,
          "x": 0,
          "y": 213
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 5,
        "points": false,
        "renderer": "flot",
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
          "y": 220
        },
        "id": 7,
        "panels": [],
        "repeat": "deployment",
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "fontSize": "100%",
        "gridPos": {
          "h": 7,
          "w": 24,
          "x": 0,
          "y": 221
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
          "y": 228
        },
        "id": 15,
        "panels": [],
        "repeat": "deployment",
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "gridPos": {
          "h": 6,
          "w": 24,
          "x": 0,
          "y": 229
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
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
          "defaults": {
            "custom": {}
          },
          "overrides": []
        },
        "fill": 1,
        "gridPos": {
          "h": 6,
          "w": 24,
          "x": 0,
          "y": 235
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
        "percentage": false,
        "pluginVersion": "7.1.1",
        "pointradius": 2,
        "points": false,
        "renderer": "flot",
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
    "schemaVersion": 26,
    "style": "dark",
    "tags": [
      "3scale",
      "backend"
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
        },
        {
          "allValue": null,
          "datasource": "$datasource",
          "definition": "label_values(kube_pod_info{namespace='$namespace',pod=~'backend-.*'}, pod)",
          "hide": 0,
          "includeAll": true,
          "label": "deployment",
          "multi": false,
          "name": "deployment",
          "options": [],
          "query": "label_values(kube_pod_info{namespace='$namespace',pod=~'backend-.*'}, pod)",
          "refresh": 1,
          "regex": "/(.*)-[a-z0-9]+-[a-z0-9]+/",
          "skipUrlSync": false,
          "sort": 1,
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
    "title": "3scale Backend"
}