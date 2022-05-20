#include <bits/stdc++.h>
using namespace std;

int main() {
	int multi = -1, xx =-1, yy = -1;
	for(int x = 1; x < 1500; x++){
		int sobra = 1500 - 2*x;
		int y = sobra / 2;
		if(y <= 0 || y > x) {
			continue;
		}

		int nodos = 2*x + y + y + 2;
		int aristas = 2*y + x + 2*x*y;
		if(nodos + aristas > multi) {
			multi = nodos+aristas;
			xx = x, yy = y;
		}
	}

	cout << "x= " << xx << " y= " << yy << " multi " << multi << endl;
}
