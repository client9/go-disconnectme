package disconnectme

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// SourceJSON is the source for the Disconnect.Me listing
const SourceJSON = "https://raw.githubusercontent.com/disconnectme/disconnect-tracking-protection/master/services.json"

// Vendor lists a entry in a Disconnect.Me category
type Vendor struct {
	Name    string   // name of vendor, e.g. "Google"
	Address string   // primary website address, e.g. "https://www.google.com/"
	Domains []string // domains under management e.g. "google.com, google.ca", "gmail.com"
}

// CategoryVendorList is map of a category to a list of Vendors
//  Categories are not enumerated for future compatibility.
//  As of 2017, the categories are:
//    - Advertising, Analytics, Content, Social, Disconnect
//
type CategoryVendorList map[string][]Vendor

// ParseURL is a convience function to read the master Disconnect.Me
//  from the network
func ParseURL() (CategoryVendorList, error) {
	resp, err := http.Get(SourceJSON)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return Parse(resp.Body)
}

// Parse takes a reader of the JSON list
func Parse(r io.Reader) (CategoryVendorList, error) {
	dm := dmMaster{}
	dec := json.NewDecoder(r)
	err := dec.Decode(&dm)
	if err != nil {
		return nil, err
	}

	out := make(map[string][]Vendor)
	for category, vendors := range dm.Categories {
		for _, vendor := range vendors {
			host := vendor.Host()
			out[category] = append(out[category], Vendor{
				Name:    vendor.Name(),
				Address: host.Name(),
				Domains: host.Domains(),
			})
		}
	}

	return out, nil
}

/****
 * Everything below here is for unmarshaling the JSON format
 * and not exposed
 */

// dmDomain maps a URL to various domains
//  This should be map[string][]string
//  but JSON source has stray junk in it so it does not unmarshal cleanly
type dmDomain map[string]interface{}

func (d dmDomain) Domains() []string {
	raw, ok := d[d.Name()].([]interface{})
	if !ok {
		panic("WTF")
	}
	list := make([]string, len(raw))
	for i, val := range raw {
		list[i] = val.(string)
	}
	return list
}

// Name returns the only value in the domain map or panics
func (d dmDomain) Name() string {
	// weird case, breaks serialization
	/*
	   "ItIsATracker": {
	     "https://itisatracker.com/": [
	       "itisatracker.com"
	     ],
	     "dnt": "eff"
	   }
	*/
	if len(d) == 2 {
		delete(d, "dnt")
	}

	if len(d) != 1 {
		panic("Vendor had more than one name")
	}
	for k := range d {
		return k
	}
	panic("Domain was empty")
}

// dmVendor maps a advertiser,etc to domains it runs
type dmVendor map[string]dmDomain

// Name returns the single vendor name from map
func (v dmVendor) Name() string {
	if len(v) != 1 {
		for k := range v {
			log.Printf("PANIC: MulitName %s", k)
		}
		panic("Vendor had more than one name")
	}
	for k := range v {
		return k
	}
	panic("Vendor was empty")
}

// Host
func (v dmVendor) Host() dmDomain {
	return v[v.Name()]
}

// dmMaster is the internal struct to unmarshal Disconnect.Me json
// it is the converted to something more go-like
type dmMaster struct {
	License    string                `json:"license"`
	Categories map[string][]dmVendor `json:"categories"`
}
