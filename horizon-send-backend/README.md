## Endpoints

## POST /upload 

Request:
multipart/form-data

    data: binary
    data_header: str   b
    name: str   b
    name_header: str   b
    hashed_key: str 
    hashed_key_salt: str   b
    max_downloads
    mime: str   b
    mime_header: str   b
    turnstile: str
    expires_after: str



    formData.append('data', new Blob([encryptedFile], { type: 'application/octet-stream' }));
    formData.append('data_header', headerBase64);
    formData.append('name', encryptedFileNameBase64);
    formData.append('name_header', nameHeaderBase64);
    formData.append('hashed_key', hashedKeyString);
    formData.append('hashed_key_salt', saltB64);
    formData.append('max_downloads', showMaxDownloads ? maxDownloads.toString() : '0');
    formData.append('mime', encryptedMimeBase64);
    formData.append('mime_header', mimeHeaderB64);
    formData.append('turnstile', turnstileToken);
    formData.append('expires_after', expiresAfter.toString());

Response:
    success: true;
    message: string;
    data: {
        fileUUID: string;
        deleteKey: string;
    }

## GET download/${uploadUuid}#${uploadKey}

b  B64URL


## GET /auth/check

-I Credentials: include
Access-Control-Allow-Credentials: true


Response:
    success: boolean;
    message: string;
    data: {
        supporter: boolean;
    };


## DELETE file/${uploadDeleteKey}



## GET /argon_salt/<fileUUID>   # using this get the hashed key
Response:
    success: boolean;
    message: string;
    data: string;   # salt value


## GET /meta/<fileUUID>/?hashed_key=""
Response:
    encrypted_name: string;
    encrypted_name_header: string;
    data_header: string;
    max_downloads: number;
    expires_at: number;
    encrypted_mime: string;
    encrypted_mime_header: string;
    size: number;
    chunk_uploaded: boolean;


## GET /download/<fileUUID>/?hashed_key=""

Response:
    success: boolean;
    message: string;
    data: string;


## GET BRIDGE_URL/?uuid=<fileUUID>&auth=<data_from_download_request>

Request:


Response:
    binary data
