# giftwrap Documentation Style Guide

This guide covers writing standards for giftwrap documentation. Follow it when adding or editing any page.

## Audience

Engineers. Readers are familiar with Go, build tooling, and the command line. Write for the person who wants the answer and will move on.

## Tone

Direct and declarative. State what things do, not how impressive they are.

**Avoid:**
> giftwrap is a powerful tool that helps you build and package your Go applications for multiple platforms.

**Prefer:**
> giftwrap cross-compiles Go applications and packages the output into release archives.

No filler. No marketing language. No "simply" or "easily."

## Structure

Each page stands alone. Do not assume the reader arrived from another page, and do not require them to read one. Include context at the top of each page.

**Hierarchy:**

- One H1 per page (the title, from front matter)
- H2 for major sections
- H3 for sub-concepts within a section
- Avoid H4 and deeper

Keep sections short. If a section runs more than three paragraphs, split it.

## Examples

Every concept gets a worked example. An example is a concrete YAML snippet, command invocation, or output block — not a hypothetical or a description.

**Weak:**
> The `exec` field runs a command before or after the build.

**Strong:**
> The `exec` field runs a command before or after the build.
>
> ```yaml
> exec:
>   pre: go generate ./...
>   post: echo "built $GOOS/$GOARCH"
> ```

Show the actual thing. Readers learn faster from a concrete example than from a description of what an example would look like.

## Formatting

- Use fenced code blocks for all YAML, shell commands, and file paths
- Use tables for structured data (auto-set variables, format options, etc.)
- Use bold for key terms on first introduction
- Do not use italics for emphasis; restructure the sentence instead

## Field documentation format (wrapfile reference)

Follow the GitHub Actions workflow syntax style:

- H2 heading for each top-level field, name in backticks: `## \`fieldname\``
- H3 heading for nested fields, using dot-notation: `### \`parent.<name>.child\``
- Required or optional stated in the first sentence of the description
- Type and default inline in prose, in backticks
- One minimal YAML example per field

```markdown
## `fieldname`

Required. One-sentence description of what the field controls.

```yaml
fieldname: example
```

### `parent.<name>.child`

Optional. Default: `value`. Description.

```yaml
parent:
  myname:
    child: example
```
```

## Language

- Use present tense: "giftwrap writes binaries to `_build/`" not "giftwrap will write"
- Use active voice: "giftwrap runs pre/post commands" not "pre/post commands are run by giftwrap"
- Spell out acronyms on first use: "operating system (OS)"
- Use `giftwrap` (lowercase) when referring to the tool

## Page checklist

Before submitting a new or updated page:

- [ ] Title clearly states what the page covers
- [ ] First paragraph explains the page's purpose without requiring prior context
- [ ] Every concept has at least one code example
- [ ] All YAML and shell snippets are in fenced code blocks
- [ ] Dates are absent from front matter
