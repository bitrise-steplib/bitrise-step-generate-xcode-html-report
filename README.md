# Generate Xcode html report

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-generate-xcode-html-report?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-generate-xcode-html-report/releases)

The Step converts xcresult summaries to html reports.

<details>
<summary>Description</summary>

Fill out.
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `verbose` | Enable logging additional information for debugging. | required | `false` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `BITRISE_HTML_REPORT_DIR` | This folder contains the generated html reports and their assets. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-generate-xcode-html-report/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-generate-xcode-html-report/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Note: this step's end-to-end tests (defined in e2e/bitrise.yml) are working with secrets which are intentionally not stored in this repo. External contributors won't be able to run those tests. Don't worry, if you open a PR with your contribution, we will help with running tests and make sure that they pass.

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
