---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by FenghuiLee.
--- DateTime: 2023/9/8 10:27
---
---
local mysql = {
    builder = {}
}

function mysql.row_field_map(map, data)
    local _data = {}
    for k, v in pairs(map) do
        if type(v) == 'string' then
            if data[v] ~= nil then
                _data[k] = data[v]
            end
        else
            if #v == 1 then
                _data[k] = v[1]
            else
                if data[v[1]] ~= nil then
                    _data[k] = v[2]
                end
            end
        end
    end
    return _data
end

function mysql.build_sql(schema, id, data, unixtimed)
    if unixtimed then
        data['updated_at'] = os.date("%Y-%m-%d %H:%M:%S", os.time())
        if id == 0 then
            data['created_at'] = data['updated_at']
        end
    end
    if id == 0 then
        return mysql.builder.insert(schema, data)
    else
        return mysql.builder.update(schema, data, id)
    end
end

function mysql.builder.insert(schema, data)
    local tCols, tVals = {}, {}
    for Key, Value in pairs(data) do
        table.insert(tCols, ('`%s`'):format(Key))
        table.insert(tVals, ("'%s'"):format(Value))
    end
    local sql = ([[
INSERT INTO `%s`.`%s`
(
    %s
)
VALUE
(
    %s
);]]):format(
        schema[1], schema[2],
        table.concat(tCols, ',\n    '),
        table.concat(tVals, ',\n    ')
    )
    print(sql)
    return sql
end

function mysql.builder.update(schema, data, id)
    local tColsVals = {}
    for Key, Value in pairs(data) do
        table.insert(tColsVals, ("`%s` = '%s'"):format(Key, Value))
    end
    local sql = ([[
UPDATE `%s`.`%s`
SET
    %s
WHERE `id` = %d;]]):format(
        schema[1], schema[2],
        table.concat(tColsVals, ',\n    '),
        id
    )
    print(sql)
    return sql
end

return mysql