# editor support

## vscode

not planned (i don't use vscode). if you are willing to make a plugin, feel free to shoot a pr

## vim

just as vscode

## neovim

using lazy:

```lua
return {
  {
    "marzeq/tree-sitter-mconf",
    dependencies = { "nvim-treesitter/nvim-treesitter" },
    config = function()
      local parser_config = require("nvim-treesitter.parsers").get_parser_configs()
      parser_config.mconf = {
        install_info = {
          -- this is a bit hacky, but it's the best way to avoid duplicating the download
          url = "~/.local/share/nvim/lazy/tree-sitter-mconf",
          files = { "src/parser.c" },
        },
      }

      vim.filetype.add({
        pattern = { [".*%.mconf"] = "mconf" },
      })

      vim.api.nvim_create_autocmd("FileType", {
        pattern = "mconf",
        callback = function()
          vim.api.nvim_command("set commentstring=#\\ %s")
        end,
      })
    end,
  },
}
```

i don't know how to make it work for a different plugin manager, so you're on your own on that one (but it should be easy to adapt)

## emacs

lisp is strictly forbidden from touching my filesystem, so no support for emacs.
just kidding, i don't use emacs, but just as with the other ones, you're welcome to make a mconf-mode
