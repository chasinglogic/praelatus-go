// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package token

import "testing"

func TestLookupIdent(t *testing.T) {
	if LookupIdent("AND") != AND {
		t.Errorf("Expected AND Got %s", LookupIdent("AND"))
	}

	if LookupIdent("and") != AND {
		t.Errorf("Expected AND Got %s", LookupIdent("and"))
	}

	if LookupIdent("OR") != OR {
		t.Errorf("Expected OR Got %s", LookupIdent("OR"))
	}

	if LookupIdent("or") != OR {
		t.Errorf("Expected OR Got %s", LookupIdent("or"))
	}

	if LookupIdent("WORD") != IDENT {
		t.Errorf("Expected IDENT Got %s", LookupIdent("WORD"))
	}
}
