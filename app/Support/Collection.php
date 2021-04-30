<?php

namespace App\Support;

use ArrayAccess;
use ArrayIterator;
use IteratorAggregate;

class Collection implements ArrayAccess, IteratorAggregate
{
    protected $data = [];
    protected $_position = 0;

    public function __construct(array $array)
    {
        $this->data = $array;
    }

    public function getIterator()
    {
        return new ArrayIterator($this->data);
    }

    public function first()
    {
        if ($this->count() === 0) {
            return null;
        }
        // https://stackoverflow.com/a/3771228
        return array_values($this->data)[0];
    }

    public function count()
    {
        return count($this->data);
    }

    public function shift()
    {
        return array_shift($this->data);
    }

    public function push($item)
    {
        return array_push($this->data, $item);
    }

    public function pop()
    {
        return array_pop($this->data);
    }

    public function filter(callable $callback)
    {
        return new Collection(array_filter($this->data, $callback));
    }

    public function map(callable $callback)
    {
        return new Collection(array_map($callback, $this->data));
    }

    public function get($key)
    {
        return $this->offsetGet($key);
    }

    public function where($key, $value) 
    {
        return $this->filter(function($item) use ($key, $value) {
            if(is_object($item)) {
                if(isset($item->$key)) {
                    return $item->$key === $value;
                }
            }
            if(is_array($item)) {
                if(isset($item[$key])) {
                    return $item[$key] === $value;
                }
            }
            return false;
        });
    }

    public function toJson()
    {
        return json_encode($this->data);
    }

    public function toArray()
    {
        return $this->data;
    }

    public function offsetExists($offset)
    {
        return isset($this->data[$offset]);
    }

    public function offsetGet($offset)
    {
        return isset($this->data[$offset]) ? $this->data[$offset] : null;
    }

    public function offsetUnset($offset)
    {
        unset($this->data[$offset]);
    }

    public function offsetSet($offset, $value)
    {
        if (is_null($offset)) {
            $this->data[] = $value;
        } else {
            $this->data[$offset] = $value;
        }
    }

    public function current()
    {
        echo __FUNCTION__;
        return $this->offsetGet($this->_position);
    }

    public function key()
    {
        echo __FUNCTION__;

        return $this->_position;
    }

    public function next()
    {
        echo __FUNCTION__;

        $this->_position++;
    }

    public function rewind()
    {
        echo __FUNCTION__;

        $this->_position = 0;
    }

    public function valid()
    {
        echo __FUNCTION__;

        return $this->offsetExists($this->_position);
    }
}
