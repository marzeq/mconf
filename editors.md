# editor support

## vscode

not planned (i don't use vscode). if you are willing to make a plugin, feel free to shoot a pr

## vim

just as vscode

## neovim 0.12

install and set-up nvim-treesitter and add this snipped to init hook

```lua
vim.filetype.add({
  pattern = { [".*%.mconf"] = "mconf" },
})

vim.api.nvim_create_autocmd("User", { pattern = "TSUpdate",
callback = function()
  require("nvim-treesitter.parsers").mconf = {
    install_info = {
      url = "https://github.com/marzeq/tree-sitter-mconf",
      revision = "f1422fe2c06c6e7f7b7ba3b48bb26364aef5fec7",
      queries = "queries/mconf",
    },
    tier = 2,
  }
end})

vim.api.nvim_create_autocmd("FileType", {
  pattern = "mconf",
  callback = function()
    vim.bo.commentstring = "# %s"
  end,
})
```

i don't know how to make it work for a different plugin manager, so you're on your own on that one (but it should be easy to adapt)

## emacs

lisp is strictly forbidden from touching my filesystem, so no support for emacs.
just kidding, i don't use emacs, but just as with the other ones, you're welcome to make a mconf-mode
