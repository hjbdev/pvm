<?php

namespace App\Support;

class ExeInfo
{
    public static function getFileVersion($fileName)
    {
        $key = "P\x00r\x00o\x00d\x00u\x00c\x00t\x00V\x00e\x00r\x00s\x00i\x00o\x00n\x00\x00\x00";
        $fptr = fopen($fileName, "rb");
        $data = "";
        while (!feof($fptr))
        {
           $data .= fread($fptr, 65536);
           if (strpos($data, $key)!==FALSE)
              break;
           $data = substr($data, strlen($data)-strlen($key));
        }
        fclose($fptr);
        if (strpos($data, $key)===FALSE)
           return "";
        $pos = strpos($data, $key)+strlen($key);
        $version = "";
        for ($i=$pos; $data[$i]!="\x00"; $i+=2)
           $version .= $data[$i];
        return $version;
    }
}
