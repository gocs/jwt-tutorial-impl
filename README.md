# jwt-tutorial-impl

i tried jwt to learn how oauth also work

## reference

Every screenshot of the file is every step/commit in this repo to learn easily.

https://auth0.com/blog/authentication-in-golang/

## run

```
$ AUTH0SECRET="" AUTH0AUDIENCE="" AUTH0DOMAIN="" go run .
```

fill the empty string above based on the values below.
<hr>

**AUTH0SECRET**

![APIs > above RBAC settings](./assets/press-copy-secret-key.png)
<hr>

**AUTH0AUDIENCE**

also called as *AUTH0_API_AUDIENCE*

![if there's an error in the domain](./assets/audience-location.png)
<hr>

**AUTH0_CALLBACK_URL**

![Route with login](./assets/localhost-callback-url.png)
<hr>

**AUTH0DOMAIN**

also called as *AUTH0_DOMAIN*

![application, locations, and its type](./assets/domain-location.png)
<hr>

**AUTH0_CLIENT_ID**

![application, locations, and its type](./assets/client-id-location.png)
<hr>

## license

MIT
