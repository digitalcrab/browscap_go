<?php

ini_set('memory_limit', '1g');

$iniFile = dirname(__FILE__) . DIRECTORY_SEPARATOR . 'test-data' . DIRECTORY_SEPARATOR . 'full_php_browscap.ini';
$ini = parse_ini_file($iniFile, true, INI_SCANNER_RAW);

echo sizeof($ini) . PHP_EOL;
