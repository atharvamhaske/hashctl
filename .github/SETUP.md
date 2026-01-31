# GitHub Actions Setup Guide

## Automatic Token (Recommended - No Setup Needed!)

GitHub Actions automatically provides a `GITHUB_TOKEN` for each workflow run. **You don't need to do anything** - the workflow will work out of the box!

The token is automatically available as `secrets.GITHUB_TOKEN` and has permissions to:
- Create releases
- Upload release assets
- Write to the repository

## How It Works

1. **Push a tag:**
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. **GitHub Actions automatically:**
   - Detects the tag push
   - Runs the workflow
   - Uses the automatic `GITHUB_TOKEN`
   - Builds all binaries
   - Creates a GitHub Release
   - Uploads all binaries

**That's it! No token setup required.**

## Using a Custom Token (Optional)

If you want to use a custom Personal Access Token (PAT) instead:

### Step 1: Create a Personal Access Token

1. Go to GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Give it a name: `hashctl-releases`
4. Select scopes:
   - ✅ `repo` (Full control of private repositories)
   - ✅ `write:packages` (if needed)
5. Click "Generate token"
6. **Copy the token immediately** (you won't see it again!)

### Step 2: Add Token as Secret

1. Go to your repository on GitHub
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Name: `GH_TOKEN` (or any name you prefer)
5. Value: Paste your token
6. Click **Add secret**

### Step 3: Update Workflow (if using custom token)

If you want to use a custom token, update `.github/workflows/release.yml`:

```yaml
env:
  GITHUB_TOKEN: ${{ secrets.GH_TOKEN }}  # Change to your secret name
```

**Note:** This is optional! The default `secrets.GITHUB_TOKEN` works perfectly.

## Troubleshooting

### Workflow fails with "permission denied"

1. Check that `permissions: contents: write` is in the workflow (it is!)
2. Make sure you're pushing the tag from a branch that has Actions enabled
3. Check the Actions tab for detailed error messages

### Release not created

1. Check the Actions tab → Your workflow run → See what failed
2. Make sure the tag format is correct: `v0.1.0`, `v1.0.0`, etc.
3. Verify the workflow file is in `.github/workflows/release.yml`

### Token permissions error

If using a custom token, make sure it has `repo` scope with write permissions.

## Testing the Workflow

1. **Test locally first:**
   ```bash
   ./build.sh v0.1.0
   ```

2. **Create a test tag:**
   ```bash
   git tag v0.1.0-test
   git push origin v0.1.0-test
   ```

3. **Check Actions tab:**
   - Go to your repo → Actions tab
   - You should see the workflow running
   - Wait for it to complete
   - Check the Releases page for the new release

4. **Delete test tag:**
   ```bash
   git tag -d v0.1.0-test
   git push origin :refs/tags/v0.1.0-test
   ```

## Summary

✅ **Default setup:** Just push a tag, everything works automatically!  
✅ **No token needed:** `GITHUB_TOKEN` is provided automatically  
✅ **Permissions set:** The workflow already has `contents: write` permission  

You're all set! Just push a tag and watch the magic happen. 🚀


