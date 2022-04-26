#!/bin/sh

swagger generate server -q -f web/swagger.yaml -P models.UserID

swagger generate server -q -f web/swagger.yaml -P models.UserID -A mindwell-images -s restapi_images \
 -O PutMeAvatar -O PutMeCover -O PostImages -O GetImagesID -O DeleteImagesID \
