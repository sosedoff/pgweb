package data

import (
	"fmt"
	"io/ioutil"
	"strings"
	"os"
	"path"
	"path/filepath"
)

// bindata_read reads the given file from disk. It returns an error on failure.
func bindata_read(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

// static_css_app_css reads file data from disk. It returns an error on failure.
func static_css_app_css() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/css/app.css"
	name := "static/css/app.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_bootstrap_css reads file data from disk. It returns an error on failure.
func static_css_bootstrap_css() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/css/bootstrap.css"
	name := "static/css/bootstrap.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_css_font_awesome_css reads file data from disk. It returns an error on failure.
func static_css_font_awesome_css() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/css/font-awesome.css"
	name := "static/css/font-awesome.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_fonts_fontawesome_otf reads file data from disk. It returns an error on failure.
func static_fonts_fontawesome_otf() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/fonts/FontAwesome.otf"
	name := "static/fonts/FontAwesome.otf"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_fonts_fontawesome_webfont_eot reads file data from disk. It returns an error on failure.
func static_fonts_fontawesome_webfont_eot() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/fonts/fontawesome-webfont.eot"
	name := "static/fonts/fontawesome-webfont.eot"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_fonts_fontawesome_webfont_svg reads file data from disk. It returns an error on failure.
func static_fonts_fontawesome_webfont_svg() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/fonts/fontawesome-webfont.svg"
	name := "static/fonts/fontawesome-webfont.svg"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_fonts_fontawesome_webfont_ttf reads file data from disk. It returns an error on failure.
func static_fonts_fontawesome_webfont_ttf() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/fonts/fontawesome-webfont.ttf"
	name := "static/fonts/fontawesome-webfont.ttf"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_fonts_fontawesome_webfont_woff reads file data from disk. It returns an error on failure.
func static_fonts_fontawesome_webfont_woff() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/fonts/fontawesome-webfont.woff"
	name := "static/fonts/fontawesome-webfont.woff"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_icon_ico reads file data from disk. It returns an error on failure.
func static_img_icon_ico() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/img/icon.ico"
	name := "static/img/icon.ico"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_img_icon_png reads file data from disk. It returns an error on failure.
func static_img_icon_png() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/img/icon.png"
	name := "static/img/icon.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_index_html reads file data from disk. It returns an error on failure.
func static_index_html() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/index.html"
	name := "static/index.html"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_ace_pgsql_js reads file data from disk. It returns an error on failure.
func static_js_ace_pgsql_js() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/js/ace-pgsql.js"
	name := "static/js/ace-pgsql.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_ace_js reads file data from disk. It returns an error on failure.
func static_js_ace_js() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/js/ace.js"
	name := "static/js/ace.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_app_js reads file data from disk. It returns an error on failure.
func static_js_app_js() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/js/app.js"
	name := "static/js/app.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// static_js_jquery_js reads file data from disk. It returns an error on failure.
func static_js_jquery_js() (*asset, error) {
	path := "/Users/sosedoff/go/src/github.com/sosedoff/pgweb/static/js/jquery.js"
	name := "static/js/jquery.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"static/css/app.css": static_css_app_css,
	"static/css/bootstrap.css": static_css_bootstrap_css,
	"static/css/font-awesome.css": static_css_font_awesome_css,
	"static/fonts/FontAwesome.otf": static_fonts_fontawesome_otf,
	"static/fonts/fontawesome-webfont.eot": static_fonts_fontawesome_webfont_eot,
	"static/fonts/fontawesome-webfont.svg": static_fonts_fontawesome_webfont_svg,
	"static/fonts/fontawesome-webfont.ttf": static_fonts_fontawesome_webfont_ttf,
	"static/fonts/fontawesome-webfont.woff": static_fonts_fontawesome_webfont_woff,
	"static/img/icon.ico": static_img_icon_ico,
	"static/img/icon.png": static_img_icon_png,
	"static/index.html": static_index_html,
	"static/js/ace-pgsql.js": static_js_ace_pgsql_js,
	"static/js/ace.js": static_js_ace_js,
	"static/js/app.js": static_js_app_js,
	"static/js/jquery.js": static_js_jquery_js,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"static": &_bintree_t{nil, map[string]*_bintree_t{
		"css": &_bintree_t{nil, map[string]*_bintree_t{
			"app.css": &_bintree_t{static_css_app_css, map[string]*_bintree_t{
			}},
			"bootstrap.css": &_bintree_t{static_css_bootstrap_css, map[string]*_bintree_t{
			}},
			"font-awesome.css": &_bintree_t{static_css_font_awesome_css, map[string]*_bintree_t{
			}},
		}},
		"fonts": &_bintree_t{nil, map[string]*_bintree_t{
			"FontAwesome.otf": &_bintree_t{static_fonts_fontawesome_otf, map[string]*_bintree_t{
			}},
			"fontawesome-webfont.eot": &_bintree_t{static_fonts_fontawesome_webfont_eot, map[string]*_bintree_t{
			}},
			"fontawesome-webfont.svg": &_bintree_t{static_fonts_fontawesome_webfont_svg, map[string]*_bintree_t{
			}},
			"fontawesome-webfont.ttf": &_bintree_t{static_fonts_fontawesome_webfont_ttf, map[string]*_bintree_t{
			}},
			"fontawesome-webfont.woff": &_bintree_t{static_fonts_fontawesome_webfont_woff, map[string]*_bintree_t{
			}},
		}},
		"img": &_bintree_t{nil, map[string]*_bintree_t{
			"icon.ico": &_bintree_t{static_img_icon_ico, map[string]*_bintree_t{
			}},
			"icon.png": &_bintree_t{static_img_icon_png, map[string]*_bintree_t{
			}},
		}},
		"index.html": &_bintree_t{static_index_html, map[string]*_bintree_t{
		}},
		"js": &_bintree_t{nil, map[string]*_bintree_t{
			"ace-pgsql.js": &_bintree_t{static_js_ace_pgsql_js, map[string]*_bintree_t{
			}},
			"ace.js": &_bintree_t{static_js_ace_js, map[string]*_bintree_t{
			}},
			"app.js": &_bintree_t{static_js_app_js, map[string]*_bintree_t{
			}},
			"jquery.js": &_bintree_t{static_js_jquery_js, map[string]*_bintree_t{
			}},
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

