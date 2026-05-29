# Phase 7: Release Automation & Publishing — Manual Steps

## Step 1: GitHub Releases (runs on tag push)

Already configured. Just tag and push:
```bash
git tag v1.0.0
git push origin v1.0.0
```

The `.github/workflows/release.yml` triggers GoReleaser which builds binaries and uploads to GitHub Releases.

---

## Step 2: Homebrew

### 2a. Create the tap repository
1. Go to https://github.com/new
2. Repository name: **homebrew-tap**
3. Description: "Homebrew tap for on-a-meet"
4. Public, no template, no README
5. Click **Create repository**

### 2b. Generate a GitHub PAT
1. Go to https://github.com/settings/tokens
2. Generate new token (classic) with **repo** scope
3. Copy the token

### 2c. Add secret to on-a-meet repo
1. Go to: https://github.com/sergiocarracedo/on-a-meet/settings/secrets/actions
2. Add `HOMEBREW_TAP_GITHUB_TOKEN` with the PAT from step 2b

### 2d. Update `.goreleaser.yaml`
Uncomment the `brews` section and replace `<YOUR_USER>` with your GitHub username.

Then tag and push — GoReleaser will auto-push the formula to your tap:
```bash
git tag v1.0.0
git push origin v1.0.0
```

### 2e. Install via Homebrew
```bash
brew tap sergiocarracedo/tap
brew install on-a-meet
```

---

## Step 3: npm

### 3a. Create npm account
1. Go to https://www.npmjs.com/signup
2. Create account (or login if you have one)

### 3b. Generate npm token
1. Go to https://www.npmjs.com/settings/tokens
2. Create **Publish** token (classic)
3. Copy the token

### 3c. Add secret to on-a-meet repo
1. Go to: https://github.com/sergiocarracedo/on-a-meet/settings/secrets/actions
2. Add `NPM_TOKEN` with the npm token

### 3d. Publish the npm package

On the next tag push, `publish-npm` job in `release.yml` will publish automatically.

Alternatively, publish manually:
```bash
cd packages/npm
npm publish --access public
```

### 3e. Install via npm
```bash
npm i -g on-a-meet
```

---

## Step 4: Full release flow

```bash
# Cut a release
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions handles the rest:
# 1. GoReleaser builds binaries + pushes Homebrew formula
# 2. npm publish job publishes the wrapper package
```

Result:
- **GitHub Release** at `https://github.com/sergiocarracedo/on-a-meet/releases/tag/v1.0.0`
- **Homebrew**: `brew tap sergiocarracedo/tap && brew install on-a-meet`
- **npm**: `npm i -g on-a-meet`
