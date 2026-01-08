#!/usr/bin/env python3

import time
import random
import requests
from datetime import datetime

PUSHGATEWAY_URL = "http://localhost:9091/metrics/job/http_requests/instance/api-server"

def push_single_metric():
    """推送单个指标到 Pushgateway"""
    print("[" + str(datetime.now()) + "] 开始推送指标到 Pushgateway...")

    # 生成一个简单的指标
    metric_data = 'http_requests_total{job="api",service="user-service",method="GET",status="200"} 500.00'

    print("推送的指标数据:")
    print(metric_data)
    print()

    try:
        response = requests.post(
            PUSHGATEWAY_URL,
            data=metric_data.encode('utf-8'),
            headers={"Content-Type": "text/plain; charset=utf-8"}
        )

        if response.status_code == 200:
            print("[" + str(datetime.now()) + "] ✅ 指标推送成功！")
        else:
            print("[" + str(datetime.now()) + "] ❌ 指标推送失败！状态码: " + str(response.status_code))
            print("响应: " + response.text)

    except Exception as e:
        print("[" + str(datetime.now()) + "] ❌ 推送指标时发生错误: " + str(e))

def main():
    print("=" * 60)
    print("Prometheus 指标推送脚本")
    print("=" * 60)
    print("Pushgateway URL: " + PUSHGATEWAY_URL)
    print("按 Ctrl+C 停止脚本")
    print("=" * 60)
    print()

    push_count = 0
    try:
        while True:
            push_count += 1
            print("\n--- 第 " + str(push_count) + " 次推送 ---")
            push_single_metric()

            # 每 10 秒推送一次
            time.sleep(10)

    except KeyboardInterrupt:
        print("\n\n[" + str(datetime.now()) + "] 停止推送")
        print("总共推送了 " + str(push_count) + " 次")

if __name__ == "__main__":
    main()
