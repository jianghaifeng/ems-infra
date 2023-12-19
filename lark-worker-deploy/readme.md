kubectl create secret generic lark-worker-secret \
--from-literal=APP_ID=<app_id> \
--from-literal=APP_SECRET=<app_secret> \
--from-literal=TEMPLATE_APP_ID=<template_id> \
--from-literal=TEMPLATE_FOLDER_ID=<folder_id> \
-n ems-prod