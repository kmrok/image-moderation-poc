# image-moderation-poc

## Overview

image-moderation-poc is a PoC repository for image moderation.
It uses AWS Rekognition API and GCP Cloud Vision API for technical validation.

## Run

You will need AWS and GCP credentials in order to run the application.
The GCP credentials should be placed in `.gcloud/credentials.json` and the AWS credentials should be passed as arguments to the makefile command as follows.

```sh
$ make run AWS_ACCESS_KEY_ID=************** AWS_SECRET_ACCESS_KEY=**************
AWS / Picnic -> 1.698094 seconds
AWS / Swimwear -> 1.835012 seconds
GCP(url) / Picnic -> 3.532163 seconds
GCP(url) / Swimwear -> 2.976572 seconds
GCP(download) / Picnic -> 1.429086 seconds
GCP(download) / Swimwear -> 1.513855 seconds
```
