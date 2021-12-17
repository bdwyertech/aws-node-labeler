#!/bin/bash

helm lint aws-node-labeler
helm package aws-node-labeler
helm repo index .
