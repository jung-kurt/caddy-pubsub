-- This script reads a single page of Markdown-like documentation and generates
-- output in three different forms: gofmt, git-flavored markdown, and
-- Caddy-style markdown.

local gsub, match, len, find, concat, insert = 
string.gsub, string.match, string.len, string.find, table.concat, table.insert

local function write(filestr, str)
  local f = io.open(filestr, 'w+')
  if f then
    f:write(str)
    f:close()
  end
end

local function codeblock(tbl, mode) 
  local newtbl = {}
  local incode = false
  local pos1, pos2, prefix, syntax
  for j, str in ipairs(tbl) do
    prefix, syntax = match(str, '^(```)(%a*)')
    if prefix and len(syntax) > 0 then
      incode = true
      if mode == 'r' then
        insert(newtbl, str)
      end
    elseif prefix then
      incode = false
      if mode == 'r' then
        insert(newtbl, str)
      end
    else
      if incode and mode ~= 'r' then
        str = '\t' .. str
      end
      insert(newtbl, str)
    end
  end
  return newtbl
end

local function caddywrite(tbl)
  tbl = codeblock(tbl, 'c')
  local str = concat(tbl, '\n')
  write('pubsub.md', str)
end

local function readmewrite(tbl)
  tbl = codeblock(tbl, 'r')
  local str = concat(tbl, '\n')
  str = gsub(str, '\n%> ', '\n')
  -- str = gsub(str, '%b<>', '')
  write('../README.md', str)
end

local function godocwrite(tbl)
  tbl = codeblock(tbl, 'g')
  for j, str in ipairs(tbl) do
    str = gsub(str, '^#+ *', '')
    tbl[j] = gsub(str, '^* ', '\n• ')
  end
  local str = concat(tbl, '\n')
  str = gsub(str, '\n\n\n+', '\n\n')
  str = gsub(str, '\n%> ', '\n')
  str = gsub(str, '`', '')
  str = gsub(str, '/%*', '\x01')
  str = gsub(str, '%*', '')
  str = gsub(str, '\x01', '\x2f*')
  -- str = gsub(str, '%b<>', '')
  -- replace [foo][bar] with foo
  str = gsub(str, '%[(%C-)%]%[%C-%]', '%1')
  str = '/*\n' .. str .. '\n*/\npackage pubsub\n'
  write('../doc.go', str)
end

local godoc, caddy, readme = {}, {}, {}
local modeg, modec, moder

for str in io.lines('doc.txt') do
  local mode = string.match(str, '^~(%a*)~$')
  if mode then
    modeg = find(mode, 'g') ~= nil
    moder = find(mode, 'r') ~= nil
    modec = find(mode, 'c') ~= nil
  else
    if modeg then
      insert(godoc, str)
    end
    if modec then
      insert(caddy, str)
    end
    if moder then
      insert(readme, str)
    end
  end
end

caddywrite(caddy)
godocwrite(godoc)
readmewrite(readme)
