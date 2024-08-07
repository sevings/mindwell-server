#!/bin/sh

swagger generate server -q -f web/swagger.yaml -P models.UserID

swagger generate server -q -f web/swagger.yaml -P models.UserID -A mindwell-images -s restapi_images \
 -O PutMeAvatar -O PutMeCover \
 -O PutThemesNameAvatar -O PutThemesNameCover \
 -O PostImages -O GetImagesFind -O GetImagesID -O DeleteImagesID \
