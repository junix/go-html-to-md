#!/usr/bin/env python3

from __future__ import annotations

import argparse
import sys

import requests


def read_url(url):
    try:
        response = requests.get(url)
        response.raise_for_status()  # Raise an exception for bad status codes
        return response.text
    except requests.RequestException as e:
        print(f"Error fetching URL: {e}", file=sys.stderr)
        sys.exit(1)


def main():
    parser = argparse.ArgumentParser(description="Read and print HTML content from a URL")
    parser.add_argument("-u", "--url", required=True, help="URL to fetch")

    args = parser.parse_args()

    html_content = read_url(args.url)
    print(html_content)


if __name__ == "__main__":
    main()
