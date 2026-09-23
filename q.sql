SELECT TABLE_NAME, COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA='ngibarbe'
  AND ((TABLE_NAME='provinsi' AND COLUMN_NAME IN ('id','kode'))
    OR (TABLE_NAME='kabupatenkota' AND COLUMN_NAME IN ('id','kode','provinsiId'))
    OR (TABLE_NAME='kecamatan' AND COLUMN_NAME IN ('id','kode','kabupatenKotaId'))
    OR (TABLE_NAME='desa' AND COLUMN_NAME IN ('id','kode','kecamatanId')))
ORDER BY TABLE_NAME, COLUMN_NAME;
