# Spin oauth login

This is a simple Spin component for adding oauth login functionality to your Spin app

## Using the oauth login template

First, configure the template:

```bash
$ spin templates install --git https://github.com/rajatjindal/oauth-login-spin --update
Copying remote template source
Installing template oauth-login...
Installed 1 template(s)

+-------------------------------------------+
| Name          Description                 |
+===========================================+
| oauth-login   Add oauth login to your app |
+-------------------------------------------+
```

Then, add this component to your existing app:

```bash
$ spin add oauth-login
```

This will create the following component in your `spin.toml`:

```toml
[variables]
auth_provider = { default = "github" }
scopes = { default = "openid,profile,email" }
client_id = { required = true, secret = true }
client_secret = { required = true, secret = true }

[[component]]
id = "oauth-login"
source = { url = "https://github.com/rajatjindal/oauth-login-spin/releases/download/v0.0.8/login.wasm", digest = "sha256:d7697dcf34989c444c164a1428cc6ab8ba02a4c113d3d344492ca54b9e4ef27d" }
allowed_http_hosts = ["https://github.com"]
key_value_stores = ["default"]
[component.trigger]
route = "/internal/login/..."
[component.config]
auth_provider = "{{ auth_provider }}"
scopes = "{{ scopes }}"
client_id = "{{ client_id }}"
client_secret = "{{ client_secret }}"
```

You can now provide client_id and client_secret for your oauth app using `--variable client_id=<client_id> --variable client_secret=<client_secret>` during deploy or you can use GitHub action as follows:

```yaml
  - name: Deploy
    uses: fermyon/actions/spin/deploy@v1
    with:
      fermyon_token: ${{ secrets.FERMYON_CLOUD_TOKEN }}
      variables: |-
        client_id=${{ secrets.OAUTH_CLIENT_ID }}
        client_secret=${{ secrets.OAUTH_CLIENT_SECRET }}
```

once added to your app, you just need to navigate user to `/internal/login/start` path to get oauth process started. once the auth is successful, the component will redirect your app to `/auth/success` (configurable via config `success_url`)