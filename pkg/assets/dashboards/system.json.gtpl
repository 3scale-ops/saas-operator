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
      "id": 85,
      "panels": [],
      "title": "App",
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
      "id": 77,
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
      "stack": true,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(rate(rails_requests_total{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+'}[1m])) by (status)",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "{{`{{status}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "System requests per second (by status)",
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
          "decimals": 1,
          "format": "reqps",
          "label": "",
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
      "id": 78,
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
      "stack": true,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(rate(rails_requests_total{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+'}[1m])) by (controller)",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "{{`{{controller}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "System requests per second (by controller)",
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
          "decimals": 1,
          "format": "reqps",
          "label": "",
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
      "id": 79,
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
      "stack": true,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(rate(rails_requests_total{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+',status=~'2[0-9]*'}[1m])) by (controller)",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "{{`{{controller}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "System 2xx requests per second (by controller)",
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
          "decimals": 1,
          "format": "reqps",
          "label": "",
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
        "w": 6,
        "x": 12,
        "y": 9
      },
      "hiddenSeries": false,
      "id": 80,
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
      "percentage": true,
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
          "expr": "sum(rate(rails_requests_total{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+',status=~'4[0-9]*'}[1m])) by (controller)",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "{{`{{controller}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "System 4xx requests per second (by controller)",
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
          "decimals": 1,
          "format": "reqps",
          "label": "",
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
        "w": 6,
        "x": 18,
        "y": 9
      },
      "hiddenSeries": false,
      "id": 81,
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
          "expr": "sum(rate(rails_requests_total{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+',status=~'5[0-9]*'}[1m])) by (controller)",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "{{`{{controller}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "System 5xx requests per second (by controller)",
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
          "decimals": 1,
          "format": "reqps",
          "label": "",
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
        "h": 7,
        "w": 24,
        "x": 0,
        "y": 17
      },
      "hiddenSeries": false,
      "id": 83,
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
          "expr": "sum(rate(rails_request_duration_seconds_sum{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+'}[1m])) / sum(rate(rails_request_duration_seconds_count{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+'}[1m]))",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "Total",
          "refId": "A"
        },
        {
          "expr": "sum(rate(rails_db_runtime_seconds_sum{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+'}[1m])) / sum(rate(rails_db_runtime_seconds_count{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+'}[1m]))",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "ActiveRecord",
          "refId": "B"
        },
        {
          "expr": "sum(rate(rails_view_runtime_seconds_sum{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+'}[1m])) / sum(rate(rails_view_runtime_seconds_count{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+'}[1m]))",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "View",
          "refId": "C"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Rails average response time",
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
          "format": "s",
          "label": "",
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
        "colorScheme": "interpolateBlues",
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
        "w": 24,
        "x": 0,
        "y": 24
      },
      "heatmap": {},
      "hideZeroBuckets": false,
      "highlightCards": true,
      "id": 87,
      "legend": {
        "show": false
      },
      "links": [],
      "reverseYBuckets": false,
      "targets": [
        {
          "expr": "sum(rate(rails_request_duration_seconds_bucket{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+'}[1m])) by (le)",
          "format": "heatmap",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "{{`{{le}}`}}",
          "refId": "A"
        }
      ],
      "timeFrom": null,
      "timeShift": null,
      "title": "Requests per second heatmap (by response time bucket in seconds)",
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
      "collapsed": true,
      "datasource": null,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 32
      },
      "id": 102,
      "panels": [
        {
          "cards": {
            "cardPadding": null,
            "cardRound": null
          },
          "color": {
            "cardColor": "#5794F2",
            "colorScale": "sqrt",
            "colorScheme": "interpolateBlues",
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
            "h": 7,
            "w": 12,
            "x": 0,
            "y": 33
          },
          "heatmap": {},
          "hideZeroBuckets": false,
          "highlightCards": true,
          "id": 103,
          "interval": "30s",
          "legend": {
            "show": false
          },
          "links": [],
          "maxDataPoints": 400,
          "repeat": "controller",
          "reverseYBuckets": false,
          "targets": [
            {
              "expr": "sum(rate(rails_request_duration_seconds_bucket{namespace='$namespace',pod=~'system-app-[a-z0-9]+-[a-z0-9]+',controller=\"$controller\"}[1m])) by (le)",
              "format": "heatmap",
              "interval": "1m",
              "intervalFactor": 10,
              "legendFormat": "{{`{{le}}`}}",
              "refId": "A"
            }
          ],
          "timeFrom": null,
          "timeShift": null,
          "title": "$controller",
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
        }
      ],
      "repeat": null,
      "title": "App  controller Request per second heatmap (by response time bucket in seconds)",
      "type": "row"
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
      "id": 52,
      "panels": [],
      "title": "Sidekiq",
      "type": "row"
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
        "w": 24,
        "x": 0,
        "y": 34
      },
      "hiddenSeries": false,
      "id": 49,
      "legend": {
        "avg": false,
        "current": false,
        "hideEmpty": true,
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
          "expr": "sum(rate(sidekiq_jobs_success_total{namespace='$namespace'}[1m]))",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "Successful jobs",
          "refId": "A"
        },
        {
          "expr": "sum(rate(sidekiq_jobs_failed_total{namespace='$namespace'}[1m]))",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "Failed jobs",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Jobs processed every second",
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
        "y": 42
      },
      "hiddenSeries": false,
      "id": 93,
      "legend": {
        "avg": false,
        "current": false,
        "hideEmpty": true,
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
          "expr": "sum(rate(sidekiq_jobs_success_total{namespace='$namespace'}[1m])) by (queue)",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "Processed in {{`{{queue}}`}}",
          "refId": "A"
        },
        {
          "expr": "sum(rate(sidekiq_jobs_failed_total{namespace='$namespace'}[1m])) by (queue)",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "Failed in {{`{{queue}}`}}",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Jobs processed by queue",
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
        "x": 12,
        "y": 42
      },
      "hiddenSeries": false,
      "id": 94,
      "legend": {
        "avg": false,
        "current": false,
        "hideEmpty": true,
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
          "expr": "sum(rate(sidekiq_jobs_success_total{namespace='$namespace'}[1m])) by (worker)",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "{{`{{worker}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Jobs processed by job type",
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
        "y": 50
      },
      "hiddenSeries": false,
      "id": 96,
      "legend": {
        "avg": false,
        "current": false,
        "hideEmpty": true,
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
          "expr": "max(sidekiq_jobs_waiting_count{namespace='$namespace'}) by (queue)",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "{{`{{queue}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Queue sizes",
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
      "id": 95,
      "legend": {
        "avg": false,
        "current": false,
        "hideEmpty": true,
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
          "expr": "max(sidekiq_queue_latency{namespace='$namespace'}) by (queue)",
          "format": "time_series",
          "interval": "1m",
          "intervalFactor": 10,
          "legendFormat": "{{`{{queue}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Queues latency",
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
          "format": "s",
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
      "id": 97,
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
          "expr": "histogram_quantile(0.95, sum(rate(sidekiq_job_runtime_seconds_bucket{namespace='$namespace'}[1m])) by (le)) ",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 10,
          "legendFormat": "95% quantile of sidekiq ",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Job processing average time",
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
          "format": "s",
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
        "x": 12,
        "y": 58
      },
      "hiddenSeries": false,
      "id": 98,
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
          "expr": "histogram_quantile(0.95, sum(rate(sidekiq_job_runtime_seconds_bucket{namespace='$namespace'}[1m])) by (worker, le)) ",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 10,
          "legendFormat": "{{`{{worker}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Job processing average time by job type",
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
          "format": "s",
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
        "y": 66
      },
      "hiddenSeries": false,
      "id": 99,
      "legend": {
        "avg": false,
        "current": false,
        "hideEmpty": true,
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
          "expr": "sum(rate(sidekiq_jobs_enqueued_total{namespace='$namespace'}[1m])) by (queue)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 10,
          "legendFormat": "{{`{{queue}}`}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Sidekiq jobs pushed by queue",
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
        "y": 66
      },
      "hiddenSeries": false,
      "id": 100,
      "legend": {
        "avg": false,
        "current": false,
        "hideEmpty": true,
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
          "expr": "max(sidekiq_jobs_retry_count{namespace='$namespace'})",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 10,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Sidekiq jobs retried",
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
      "collapsed": false,
      "datasource": null,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 74
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
        "y": 75
      },
      "hideTimeOverride": true,
      "id": 64,
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
          "expr": "sum(kube_deployment_status_replicas_available{namespace='$namespace',deployment=~'$deployment'}) or kube_statefulset_status_replicas{namespace='$namespace',statefulset='$deployment'}",
          "format": "time_series",
          "interval": "",
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
        "y": 75
      },
      "hideTimeOverride": true,
      "id": 65,
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
          "expr": "sum(kube_deployment_status_replicas_unavailable{namespace='$namespace',deployment=~'$deployment'}) or (kube_statefulset_status_replicas {namespace='$namespace',statefulset='$deployment'} - kube_statefulset_status_replicas_ready {namespace='$namespace',statefulset='$deployment'})",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 1,
          "legendFormat": "",
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
        "y": 75
      },
      "hideTimeOverride": true,
      "id": 66,
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
          "expr": "count(count(container_memory_working_set_bytes{namespace='$namespace',pod=~'$deployment-.*'}) by (node))",
          "format": "time_series",
          "interval": "",
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
        "y": 75
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
          "expr": "max(sum(delta(kube_pod_container_status_restarts_total{namespace='$namespace',pod=~'$deployment-.*'}[5m]))) by (pod)",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 1,
          "legendFormat": "",
          "refId": "B"
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
        "y": 78
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
      "maxPerRow": 3,
      "nullPointMode": "null as zero",
      "percentage": false,
      "pluginVersion": "7.1.1",
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "repeat": null,
      "repeatDirection": null,
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "kube_deployment_status_replicas{namespace='$namespace',deployment=~'$deployment'} or kube_statefulset_status_replicas {namespace='$namespace',statefulset='$deployment'}",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "{{`{{deployment}}`}}-total-pods",
          "legendLink": null,
          "refId": "A",
          "step": 10
        },
        {
          "expr": "kube_deployment_status_replicas_available{namespace='$namespace',deployment=~'$deployment'} or kube_statefulset_status_replicas_ready {namespace='$namespace',statefulset='$deployment'}",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 1,
          "legendFormat": "{{`{{deployment}}`}}-avail-pods",
          "refId": "B"
        },
        {
          "expr": "kube_deployment_status_replicas_unavailable{namespace='$namespace',deployment=~'$deployment'} or (kube_statefulset_status_replicas {namespace='$namespace',statefulset='$deployment'} - kube_statefulset_status_replicas_ready {namespace='$namespace',statefulset='$deployment'})",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 1,
          "legendFormat": "{{`{{deployment}}`}}-unavail-pods",
          "refId": "C"
        },
        {
          "expr": "count(count(container_memory_working_set_bytes{namespace='$namespace',pod=~'$deployment-.*'}) by (node))",
          "format": "time_series",
          "interval": "",
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
        "y": 85
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
      "maxPerRow": 3,
      "nullPointMode": "null",
      "percentage": false,
      "pluginVersion": "7.1.1",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "repeat": null,
      "repeatDirection": null,
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(delta(kube_pod_container_status_restarts_total{namespace='$namespace',pod=~'$deployment-.*'}[5m])) by (pod)",
          "format": "time_series",
          "interval": "",
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
        "y": 91
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
        "y": 92
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
      "maxPerRow": 3,
      "nullPointMode": "null as zero",
      "percentage": false,
      "pluginVersion": "7.1.1",
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "repeat": null,
      "repeatDirection": "h",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod)",
          "format": "time_series",
          "interval": "",
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
        "y": 99
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
        "y": 100
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
      "maxPerRow": 3,
      "nullPointMode": "null as zero",
      "pageSize": null,
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "repeat": null,
      "repeatDirection": "h",
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
          "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod)",
          "format": "table",
          "instant": true,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A",
          "step": 10
        },
        {
          "expr": "sum(kube_pod_container_resource_requests{unit='core',namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod)",
          "format": "table",
          "instant": true,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "B",
          "step": 10
        },
        {
          "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod) / sum(kube_pod_container_resource_requests{unit='core',namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod)",
          "format": "table",
          "instant": true,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "C",
          "step": 10
        },
        {
          "expr": "sum(kube_pod_container_resource_limits{unit='core',namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod)",
          "format": "table",
          "instant": true,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "D",
          "step": 10
        },
        {
          "expr": "sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod) / sum(kube_pod_container_resource_limits{unit='core',namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod)",
          "format": "table",
          "instant": true,
          "interval": "",
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
        "y": 107
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
        "y": 108
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
          "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace',pod=~'$deployment-.*', container!=''}) by (pod)",
          "format": "time_series",
          "interval": "",
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
        "y": 115
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
        "y": 116
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
          "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace',pod=~'$deployment-.*', container!=''}) by (pod)",
          "format": "table",
          "instant": true,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A",
          "step": 10
        },
        {
          "expr": "sum(kube_pod_container_resource_requests{unit='byte',namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod)",
          "format": "table",
          "instant": true,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "B",
          "step": 10
        },
        {
          "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace',pod=~'$deployment-.*', container!=''}) by (pod) / sum(kube_pod_container_resource_requests{unit='byte',namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod)",
          "format": "table",
          "instant": true,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "C",
          "step": 10
        },
        {
          "expr": "sum(kube_pod_container_resource_limits{unit='byte',namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod)",
          "format": "table",
          "instant": true,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "D",
          "step": 10
        },
        {
          "expr": "sum(container_memory_working_set_bytes{namespace=~'$namespace',pod=~'$deployment-.*', container!=''}) by (pod) / sum(kube_pod_container_resource_limits{unit='byte',namespace=~'$namespace',pod=~'$deployment-.*'}) by (pod)",
          "format": "table",
          "instant": true,
          "interval": "",
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
        "y": 123
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
      "fillGradient": 0,
      "gridPos": {
        "h": 7,
        "w": 24,
        "x": 0,
        "y": 124
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
          "expr": "sum(irate(container_network_receive_bytes_total{namespace=~'$namespace',pod=~'$deployment-.*'}[5m])) by (pod)",
          "format": "time_series",
          "interval": "",
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
      "fillGradient": 0,
      "gridPos": {
        "h": 7,
        "w": 24,
        "x": 0,
        "y": 131
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
          "expr": "sum(irate(container_network_transmit_bytes_total{namespace=~'$namespace',pod=~'$deployment-.*'}[5m])) by (pod)",
          "format": "time_series",
          "interval": "",
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
          "selected": false,
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
        "allValue": "",
        "datasource": "$datasource",
        "definition": "label_values(kube_pod_info{namespace='$namespace',pod=~'system-.*'}, pod)",
        "hide": 0,
        "includeAll": true,
        "label": "deployment",
        "multi": false,
        "name": "deployment",
        "options": [],
        "query": "label_values(kube_pod_info{namespace='$namespace',pod=~'system-.*'}, pod)",
        "refresh": 1,
        "regex": "/^([a-z-]+)-([0-9]|[a-z0-9]+-[a-z0-9]+)$/",
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
        "datasource": "$datasource",
        "definition": "label_values(rails_request_duration_seconds_bucket{namespace='$namespace',pod=~'system-app-.*'}, controller)",
        "hide": 0,
        "includeAll": true,
        "label": "controller",
        "multi": true,
        "name": "controller",
        "options": [],
        "query": "label_values(rails_request_duration_seconds_bucket{namespace='$namespace',pod=~'system-app-.*'}, controller)",
        "refresh": 1,
        "regex": "",
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
  "title": "3scale System"
}