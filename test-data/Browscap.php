<?php

namespace phpbrowscap;

class Browscap
{
    public $iterations = 0;

    /**
     * Current version of the class.
     */
    const VERSION = '2.0.3';

    const CACHE_FILE_VERSION = '2.0.3';

    const UPDATE_FOPEN = 'URL-wrapper';
    const UPDATE_FSOCKOPEN = 'socket';
    const UPDATE_CURL = 'cURL';
    const UPDATE_LOCAL = 'local';

    /**
     * Options for regex patterns.
     *
     * REGEX_DELIMITER: Delimiter of all the regex patterns in the whole class.
     * REGEX_MODIFIERS: Regex modifiers.
     */
    const REGEX_DELIMITER = '@';
    const REGEX_MODIFIERS = 'i';
    const COMPRESSION_PATTERN_START = '@';
    const COMPRESSION_PATTERN_DELIMITER = '|';

    /**
     * The values to quote in the ini file
     */
    const VALUES_TO_QUOTE = 'Browser|Parent';

    const BROWSCAP_VERSION_KEY = 'GJK_Browscap_Version';

    /**
     * The headers to be sent for checking the version and requesting the file.
     */
    const REQUEST_HEADERS = "GET %s HTTP/1.0\r\nHost: %s\r\nUser-Agent: %s\r\nConnection: Close\r\n\r\n";

    /**
     * The path of the local version of the browscap.ini file from which to
     * update (to be set only if used).
     *
     * @var string
     */
    public $localFile = null;

    /**
     * The useragent to include in the requests made by the class during the
     * update process.
     *
     * @var string
     */
    public $userAgent = 'Browser Capabilities Project - PHP Browscap/%v %m';

    /**
     * Flag to enable only lowercase indexes in the result.
     * The cache has to be rebuilt in order to apply this option.
     *
     * @var bool
     */
    public $lowercase = false;

    /**
     * Flag to enable/disable silent error management.
     * In case of an error during the update process the class returns an empty
     * array/object if the update process can't take place and the browscap.ini
     * file does not exist.
     *
     * @var bool
     */
    public $silent = false;

    /**
     * Where to store the cached PHP arrays.
     *
     * @var string
     */
    public $cacheFilename = 'cache.php';

    /**
     * Where to store the downloaded ini file.
     *
     * @var string
     */
    public $iniFilename = 'browscap.ini';

    /**
     * Path to the cache directory
     *
     * @var string
     */
    public $cacheDir = null;

    /**
     * Flag to be set to true after loading the cache
     *
     * @var bool
     */
    protected $_cacheLoaded = false;

    /**
     * Where to store the value of the included PHP cache file
     *
     * @var array
     */
    protected $_userAgents = array();
    protected $_browsers = array();
    protected $_patterns = array();
    protected $_properties = array();
    protected $_source_version;

    public function __construct($cache_dir)
    {
        // has to be set to reach E_STRICT compatibility, does not affect system/app settings
        date_default_timezone_set('UTC');

        if (!isset($cache_dir)) {
            throw new Exception('You have to provide a path to read/store the browscap cache file');
        }

        $old_cache_dir = $cache_dir;
        $cache_dir     = realpath($cache_dir);

        if (false === $cache_dir) {
            throw new Exception(
                sprintf(
                    'The cache path %s is invalid. Are you sure that it exists and that you have permission to access it?',
                    $old_cache_dir
                )
            );
        }

        // Is the cache dir really the directory or is it directly the file?
        if (substr($cache_dir, -4) === '.php') {
            $this->cacheFilename = basename($cache_dir);
            $this->cacheDir      = dirname($cache_dir);
        } else {
            $this->cacheDir = $cache_dir;
        }

        $this->cacheDir .= DIRECTORY_SEPARATOR;
    }

    /**
     * @return mixed
     */
    public function getSourceVersion()
    {
        return $this->_source_version;
    }

    /**
     * XXX parse
     *
     * Gets the information about the browser by User Agent
     *
     * @param string $user_agent   the user agent string
     * @param bool   $return_array whether return an array or an object
     *
     * @throws Exception
     * @return \stdClass|array  the object containing the browsers details. Array if
     *                    $return_array is set to true.
     */
    public function getBrowser($user_agent = null, $return_array = false)
    {
        $this->iterations = 0;

        // Load the cache at the first request
        if (!$this->_cacheLoaded) {
            $cache_file = $this->cacheDir . $this->cacheFilename;
            $ini_file   = $this->cacheDir . $this->iniFilename;

            $update_cache = true;

            if (file_exists($cache_file) && file_exists($ini_file)) {
                if ($this->_loadCache($cache_file)) {
                    $update_cache = false;
                }
            }

            if ($update_cache) {
                try {
                    $this->updateCache();
                } catch (Exception $e) {
                    if (!$this->silent) {
                        throw $e;
                    }
                }

                if (!$this->_loadCache($cache_file)) {
                    throw new Exception('Cannot load this cache version - the cache format is not compatible.');
                }
            }
        }

        // Automatically detect the useragent
        if (!isset($user_agent)) {
            if (isset($_SERVER['HTTP_USER_AGENT'])) {
                $user_agent = $_SERVER['HTTP_USER_AGENT'];
            } else {
                $user_agent = '';
            }
        }

        $browser = array();
        foreach ($this->_patterns as $pattern => $pattern_data) {
            $this->iterations++;
            if (preg_match($pattern . 'i', $user_agent, $matches)) {
                if (1 == count($matches)) {
                    // standard match
                    $key = $pattern_data;

                    $simple_match = true;
                } else {
                    $pattern_data = unserialize($pattern_data);

                    // match with numeric replacements
                    array_shift($matches);

                    $match_string = self::COMPRESSION_PATTERN_START
                        . implode(self::COMPRESSION_PATTERN_DELIMITER, $matches);

                    if (!isset($pattern_data[$match_string])) {
                        // partial match - numbers are not present, but everything else is ok
                        continue;
                    }

                    $key = $pattern_data[$match_string];

                    $simple_match = false;
                }

                $browser = array(
                    $user_agent, // Original useragent
                    trim(strtolower($pattern), self::REGEX_DELIMITER),
                    $this->_pregUnQuote($pattern, $simple_match ? false : $matches)
                );

                $browser = $value = $browser + unserialize($this->_browsers[$key]);

                while (array_key_exists(3, $value)) {
                    $value = unserialize($this->_browsers[$value[3]]);
                    $browser += $value;
                }

                if (!empty($browser[3])) {
                    $browser[3] = $this->_userAgents[$browser[3]];
                }

                break;
            }
        }

        // Add the keys for each property
        $array = array();
        foreach ($browser as $key => $value) {
            if ($value === 'true') {
                $value = true;
            } elseif ($value === 'false') {
                $value = false;
            }

            $tmp_key = $this->_properties[$key];
            if ($this->lowercase) {
                $tmp_key = strtolower($this->_properties[$key]);
            }
            $array[$tmp_key] = $value;
        }

        return $return_array ? $array : (object) $array;
    }

    /**
     * XXX save
     *
     * Parses the ini file and updates the cache files
     *
     * @throws Exception
     * @return bool whether the file was correctly written to the disk
     */
    public function updateCache()
    {
        $lockfile = $this->cacheDir . 'cache.lock';

        if (file_exists($lockfile) || !touch($lockfile)) {
            throw new Exception('temporary file already exists');
        }

        $ini_path   = $this->cacheDir . $this->iniFilename;
        $cache_path = $this->cacheDir . $this->cacheFilename;

        if (version_compare(PHP_VERSION, '5.3.0', '>=')) {
            $browsers = parse_ini_file($ini_path, true, INI_SCANNER_RAW);
        } else {
            $browsers = parse_ini_file($ini_path, true);
        }

        $this->_source_version = $browsers[self::BROWSCAP_VERSION_KEY]['Version'];

        unset($browsers[self::BROWSCAP_VERSION_KEY]);
        unset($browsers['DefaultProperties']['RenderingEngine_Description']);

        $this->_properties = array_keys($browsers['DefaultProperties']);

        array_unshift(
            $this->_properties,
            'browser_name',
            'browser_name_regex',
            'browser_name_pattern',
            'Parent'
        );

        $tmp_user_agents = array_keys($browsers);

        usort($tmp_user_agents, array($this, 'compareBcStrings'));

        $user_agents_keys = array_flip($tmp_user_agents);
        $properties_keys  = array_flip($this->_properties);

        $tmp_patterns = array();

        foreach ($tmp_user_agents as $i => $user_agent) {

            if (empty($browsers[$user_agent]['Comment'])
                || false !== strpos($user_agent, '*')
                || false !== strpos($user_agent, '?')
            ) {
                $pattern = $this->_pregQuote($user_agent);

                $matches_count = preg_match_all('@\d@', $pattern, $matches);

                if (!$matches_count) {
                    $tmp_patterns[$pattern] = $i;
                } else {
                    $compressed_pattern = preg_replace('@\d@', '(\d)', $pattern);

                    if (!isset($tmp_patterns[$compressed_pattern])) {
                        $tmp_patterns[$compressed_pattern] = array('first' => $pattern);
                    }

                    $tmp_patterns[$compressed_pattern][$i] = $matches[0];
                }
            }

            if (!empty($browsers[$user_agent]['Parent'])) {
                $parent = $browsers[$user_agent]['Parent'];

                $parent_key = $user_agents_keys[$parent];

                $browsers[$user_agent]['Parent']       = $parent_key;
                $this->_userAgents[$parent_key . '.0'] = $tmp_user_agents[$parent_key];
            };

            $browser = array();
            foreach ($browsers[$user_agent] as $key => $value) {
                if (!isset($properties_keys[$key])) {
                    continue;
                }

                $key           = $properties_keys[$key];
                $browser[$key] = $value;
            }

            $this->_browsers[] = $browser;
        }

        // reducing memory usage by unsetting $tmp_user_agents
        unset($tmp_user_agents);

        foreach ($tmp_patterns as $pattern => $pattern_data) {
            if (is_int($pattern_data)) {
                $this->_patterns[$pattern] = $pattern_data;
            } elseif (2 == count($pattern_data)) {
                end($pattern_data);
                $this->_patterns[$pattern_data['first']] = key($pattern_data);
            } else {
                unset($pattern_data['first']);

                $pattern_data = $this->deduplicateCompressionPattern($pattern_data, $pattern);

                $this->_patterns[$pattern] = $pattern_data;
            }
        }

        // Get the whole PHP code
        $cache = $this->_buildCache();
        $dir   = dirname($cache_path);

        // "tempnam" did not work with VFSStream for tests
        $tmpFile = $dir . '/temp_' . md5(time() . basename($cache_path));

        // asume that all will be ok
        if (false === file_put_contents($tmpFile, $cache)) {
            // writing to the temparary file failed
            throw new Exception('wrting to temporary file failed');
        }

        if (false === rename($tmpFile, $cache_path)) {
            // renaming file failed, remove temp file
            @unlink($tmpFile);

            throw new Exception('could not rename temporary file to the cache file');
        }

        @unlink($lockfile);

        return true;
    }

    /**
     * @param string $a
     * @param string $b
     *
     * @return int
     */
    protected function compareBcStrings($a, $b)
    {
        $a_len = strlen($a);
        $b_len = strlen($b);

        if ($a_len > $b_len) {
            return -1;
        }

        if ($a_len < $b_len) {
            return 1;
        }

        $a_len = strlen(str_replace(array('*', '?'), '', $a));
        $b_len = strlen(str_replace(array('*', '?'), '', $b));

        if ($a_len > $b_len) {
            return -1;
        }

        if ($a_len < $b_len) {
            return 1;
        }

        return 0;
    }

    /**
     * That looks complicated...
     *
     * All numbers are taken out into $matches, so we check if any of those numbers are identical
     * in all the $matches and if they are we restore them to the $pattern, removing from the $matches.
     * This gives us patterns with "(\d)" only in places that differ for some matches.
     *
     * @param array  $matches
     * @param string $pattern
     *
     * @return array of $matches
     */
    protected function deduplicateCompressionPattern($matches, &$pattern)
    {
        $tmp_matches = $matches;
        $first_match = array_shift($tmp_matches);
        $differences = array();

        foreach ($tmp_matches as $some_match) {
            $differences += array_diff_assoc($first_match, $some_match);
        }

        $identical = array_diff_key($first_match, $differences);

        $prepared_matches = array();

        foreach ($matches as $i => $some_match) {
            $key = self::COMPRESSION_PATTERN_START
                . implode(self::COMPRESSION_PATTERN_DELIMITER, array_diff_assoc($some_match, $identical));

            $prepared_matches[$key] = $i;
        }

        $pattern_parts = explode('(\d)', $pattern);

        foreach ($identical as $position => $value) {
            $pattern_parts[$position + 1] = $pattern_parts[$position] . $value . $pattern_parts[$position + 1];
            unset($pattern_parts[$position]);
        }

        $pattern = implode('(\d)', $pattern_parts);

        return $prepared_matches;
    }

    /**
     * Converts browscap match patterns into preg match patterns.
     *
     * @param string $user_agent
     *
     * @return string
     */
    protected function _pregQuote($user_agent)
    {
        $pattern = preg_quote($user_agent, self::REGEX_DELIMITER);

        // the \\x replacement is a fix for "Der gro\xdfe BilderSauger 2.00u" user agent match

        return self::REGEX_DELIMITER
        . '^'
        . str_replace(array('\*', '\?', '\\x'), array('.*', '.', '\\\\x'), $pattern)
        . '$'
        . self::REGEX_DELIMITER;
    }

    /**
     * Converts preg match patterns back to browscap match patterns.
     *
     * @param string        $pattern
     * @param array|boolean $matches
     *
     * @return string
     */
    protected function _pregUnQuote($pattern, $matches)
    {
        // list of escaped characters: http://www.php.net/manual/en/function.preg-quote.php
        // to properly unescape '?' which was changed to '.', I replace '\.' (real dot) with '\?', then change '.' to '?' and then '\?' to '.'.
        $search  = array(
            '\\' . self::REGEX_DELIMITER, '\\.', '\\\\', '\\+', '\\[', '\\^', '\\]', '\\$', '\\(', '\\)', '\\{', '\\}',
            '\\=', '\\!', '\\<', '\\>', '\\|', '\\:', '\\-', '.*', '.', '\\?'
        );
        $replace = array(
            self::REGEX_DELIMITER, '\\?', '\\', '+', '[', '^', ']', '$', '(', ')', '{', '}', '=', '!', '<', '>', '|',
            ':', '-', '*', '?', '.'
        );

        $result = substr(str_replace($search, $replace, $pattern), 2, -2);

        if ($matches) {
            foreach ($matches as $one_match) {
                $num_pos = strpos($result, '(\d)');
                $result  = substr_replace($result, $one_match, $num_pos, 4);
            }
        }

        return $result;
    }

    /**
     * Loads the cache into object's properties
     *
     * @param string $cache_file
     *
     * @return boolean
     */
    protected function _loadCache($cache_file)
    {
        $cache_version  = null;
        $source_version = null;
        $browsers       = array();
        $userAgents     = array();
        $patterns       = array();
        $properties     = array();

        $this->_cacheLoaded = false;

        require $cache_file;

        if (!isset($cache_version) || $cache_version != self::CACHE_FILE_VERSION) {
            return false;
        }

        $this->_source_version = $source_version;
        $this->_browsers       = $browsers;
        $this->_userAgents     = $userAgents;
        $this->_patterns       = $patterns;
        $this->_properties     = $properties;

        $this->_cacheLoaded = true;

        return true;
    }

    /**
     * Parses the array to cache and creates the PHP string to write to disk
     *
     * @return string the PHP string to save into the cache file
     */
    protected function _buildCache()
    {
        $cacheTpl = "<?php\n\$source_version=%s;\n\$cache_version=%s;\n\$properties=%s;\n\$browsers=%s;\n\$userAgents=%s;\n\$patterns=%s;\n";

        $propertiesArray = $this->_array2string($this->_properties);
        $patternsArray   = $this->_array2string($this->_patterns);
        $userAgentsArray = $this->_array2string($this->_userAgents);
        $browsersArray   = $this->_array2string($this->_browsers);

        return sprintf(
            $cacheTpl,
            "'" . $this->_source_version . "'",
            "'" . self::CACHE_FILE_VERSION . "'",
            $propertiesArray,
            $browsersArray,
            $userAgentsArray,
            $patternsArray
        );
    }

    /**
     * Converts the given array to the PHP string which represent it.
     * This method optimizes the PHP code and the output differs form the
     * var_export one as the internal PHP function does not strip whitespace or
     * convert strings to numbers.
     *
     * @param array $array the array to parse and convert
     *
     * @return string the array parsed into a PHP string
     */
    protected function _array2string($array)
    {
        $strings = array();

        foreach ($array as $key => $value) {
            if (is_int($key)) {
                $key = '';
            } elseif (ctype_digit((string) $key) || '.0' === substr($key, -2)) {
                $key = intval($key) . '=>';
            } else {
                $key = "'" . str_replace("'", "\'", $key) . "'=>";
            }

            if (is_array($value)) {
                $value = "'" . addcslashes(serialize($value), "'") . "'";
            } elseif (ctype_digit((string) $value)) {
                $value = intval($value);
            } else {
                $value = "'" . str_replace("'", "\'", $value) . "'";
            }

            $strings[] = $key . $value;
        }

        return "array(\n" . implode(",\n", $strings) . "\n)";
    }
}

class Exception extends \Exception
{
    // nothing to do here
}
