pkgname=marco
pkgver=0.1.0
pkgrel=1
pkgdesc='Marco: a settings tui'
arch=('x86_64')
url="https://github.com/sokryptk/marco"
source=('git+https://github.com/sokryptk/marco.git')
makedepends=('go')
sha256sums=('SKIP')

build() {
    cd $pkgname
    go build -o $pkgname .
} 

package() {
    cd $pkgname
    install -Dm755 $pkgname "$pkgdir"/usr/bin/$pkgname
}
