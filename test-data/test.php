<?php
ini_set('memory_limit', '1g');
include_once "Browscap.php";

$start = microtime(true);

$cacheDir = dirname(__FILE__);

$brows = new \phpbrowscap\Browscap($cacheDir);
$brows->iniFilename = 'full_php_browscap.ini';
$res = $brows->getBrowser('Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36', true);

echo (empty($res['Browser']) ? 'false' : $res['Browser']) . PHP_EOL;
echo (microtime(true) - $start) . ' [' . $brows->iterations . ']' . PHP_EOL;