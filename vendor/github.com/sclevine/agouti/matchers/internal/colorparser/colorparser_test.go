package colorparser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers/internal/colorparser"
)

var _ = Describe("Parsing CSS Colors", func() {
	Describe("Color", func() {
		It("should String() nicely", func() {
			c := Color{R: 128, G: 255, B: 35, A: 0.3}
			Expect(c.String()).To(Equal("Color{R:128, G:255, B:35, A:0.30}"))
		})
	})

	Context("with hex notation", func() {
		Context("with #XXX hex notation", func() {
			It("should parse the color correctly", func() {
				Expect(ParseCSSColor("#000")).To(Equal(Color{0, 0, 0, 1.0}))

				Expect(ParseCSSColor("#F00")).To(Equal(Color{255, 0, 0, 1.0}))
				Expect(ParseCSSColor("#0F0")).To(Equal(Color{0, 255, 0, 1.0}))
				Expect(ParseCSSColor("#00F")).To(Equal(Color{0, 0, 255, 1.0}))
				Expect(ParseCSSColor("#FFF")).To(Equal(Color{255, 255, 255, 1.0}))

				Expect(ParseCSSColor("#a00")).To(Equal(Color{170, 0, 0, 1.0}))
				Expect(ParseCSSColor("#0a0")).To(Equal(Color{0, 170, 0, 1.0}))
				Expect(ParseCSSColor("#00a")).To(Equal(Color{0, 0, 170, 1.0}))
				Expect(ParseCSSColor("#aaa")).To(Equal(Color{170, 170, 170, 1.0}))

				Expect(ParseCSSColor("  #aaa  ")).To(Equal(Color{170, 170, 170, 1.0}))
			})
		})

		Context("with #XXXXXX hex notation", func() {
			It("should parse the color correctly", func() {
				Expect(ParseCSSColor("#000000")).To(Equal(Color{0, 0, 0, 1.0}))

				Expect(ParseCSSColor("#FF0000")).To(Equal(Color{255, 0, 0, 1.0}))
				Expect(ParseCSSColor("#00FF00")).To(Equal(Color{0, 255, 0, 1.0}))
				Expect(ParseCSSColor("#0000FF")).To(Equal(Color{0, 0, 255, 1.0}))
				Expect(ParseCSSColor("#FFFFFF")).To(Equal(Color{255, 255, 255, 1.0}))

				Expect(ParseCSSColor("#ab0000")).To(Equal(Color{171, 0, 0, 1.0}))
				Expect(ParseCSSColor("#00ab00")).To(Equal(Color{0, 171, 0, 1.0}))
				Expect(ParseCSSColor("#0000ab")).To(Equal(Color{0, 0, 171, 1.0}))
				Expect(ParseCSSColor("#abcdef")).To(Equal(Color{171, 205, 239, 1.0}))

				Expect(ParseCSSColor("  #abcdef  ")).To(Equal(Color{171, 205, 239, 1.0}))
			})
		})
	})

	Context("when passed an rgb color", func() {
		Context("with rgb(X, X, X) notation", func() {
			It("should parse the color correctly", func() {
				Expect(ParseCSSColor("rgb(0,0,0)")).To(Equal(Color{0, 0, 0, 1.0}))
				Expect(ParseCSSColor("rgb(0, 0, 0)")).To(Equal(Color{0, 0, 0, 1.0}))
				Expect(ParseCSSColor("  rgb(0, 0, 0)  ")).To(Equal(Color{0, 0, 0, 1.0}))
				Expect(ParseCSSColor("rgb(255, 128, 37)")).To(Equal(Color{255, 128, 37, 1.0}))
			})

			It("should clip colors that are out of bounds", func() {
				Expect(ParseCSSColor("rgb(-100,-100,-100)")).To(Equal(Color{0, 0, 0, 1.0}))
				Expect(ParseCSSColor("rgb(256, 257, 258)")).To(Equal(Color{255, 255, 255, 1.0}))
			})
		})

		Context("with rgb(X%, X%, X%) notation", func() {
			It("should parse the color correctly", func() {
				Expect(ParseCSSColor(`rgb(5%, 50%, 89%)`)).To(Equal(Color{13, 128, 227, 1.0}))
				Expect(ParseCSSColor(`rgb(49.9%,50%,89%)`)).To(Equal(Color{127, 128, 227, 1.0}))
			})

			It("should clip colors that are out of bounds", func() {
				Expect(ParseCSSColor(`rgb(-1%, -10%, -30%)`)).To(Equal(Color{0, 0, 0, 1.0}))
				Expect(ParseCSSColor(`rgb(100%, 101%, 102%)`)).To(Equal(Color{255, 255, 255, 1.0}))
			})
		})
	})

	Context("when passed an rgba color", func() {
		Context("with rgba(X, X, X, X) notation", func() {
			It("should parse the color correctly", func() {
				Expect(ParseCSSColor("rgba(0,0,0, 0.3)")).To(Equal(Color{0, 0, 0, 0.3}))
				Expect(ParseCSSColor("rgba(0, 0, 0, 0.7)")).To(Equal(Color{0, 0, 0, 0.7}))
				Expect(ParseCSSColor("  rgba(0, 0, 0, 1)  ")).To(Equal(Color{0, 0, 0, 1.0}))
				Expect(ParseCSSColor("rgba(255, 128, 37, 0)")).To(Equal(Color{255, 128, 37, 0.0}))
			})

			It("should clip colors that are out of bounds", func() {
				Expect(ParseCSSColor("rgba(-100,-100,-100,-0.3)")).To(Equal(Color{0, 0, 0, 0}))
				Expect(ParseCSSColor("rgba(256, 257, 258, 1.2)")).To(Equal(Color{255, 255, 255, 1.0}))
			})
		})

		Context("with rgba(X%, X%, X%, X) notation", func() {
			It("should parse the color correctly", func() {
				Expect(ParseCSSColor(`rgba(5%, 50%, 89%,0.3)`)).To(Equal(Color{13, 128, 227, 0.3}))
				Expect(ParseCSSColor(`rgba(49.9%,50%,89%,0.23)`)).To(Equal(Color{127, 128, 227, 0.23}))
			})

			It("should clip colors that are out of bounds", func() {
				Expect(ParseCSSColor(`rgba(-1%, -10%, -30%, -0.3)`)).To(Equal(Color{0, 0, 0, 0.0}))
				Expect(ParseCSSColor(`rgba(100%, 101%, 102%, 1.2)`)).To(Equal(Color{255, 255, 255, 1.0}))
			})
		})
	})

	Context("when passed an hsl color", func() {
		It("should parse the color correctly", func() {
			Expect(ParseCSSColor(`hsl(0,0%,0%)`)).To(Equal(Color{0, 0, 0, 1.0}))
			Expect(ParseCSSColor(`hsl(0,100%,50%)`)).To(Equal(Color{255, 0, 0, 1.0}))
			Expect(ParseCSSColor(`hsl(120,100%,50%)`)).To(Equal(Color{0, 255, 0, 1.0}))
			Expect(ParseCSSColor(`hsl(240,100%,50%)`)).To(Equal(Color{0, 0, 255, 1.0}))
			Expect(ParseCSSColor(`hsl(0,100%,100%)`)).To(Equal(Color{255, 255, 255, 1.0}))
			Expect(ParseCSSColor(`hsl(0,0%,100%)`)).To(Equal(Color{255, 255, 255, 1.0}))
			Expect(ParseCSSColor(`hsl(240,50%,50%)`)).To(Equal(Color{63, 63, 191, 1.0}))
			Expect(ParseCSSColor(`hsl(359,50%,50%)`)).To(Equal(Color{191, 63, 65, 1.0}))
		})

		It("should wrap the hue circle correctly", func() {
			Expect(ParseCSSColor(`hsl(0,100%,50%)`)).To(Equal(Color{255, 0, 0, 1.0}))
			Expect(ParseCSSColor(`hsl(-120,100%,50%)`)).To(Equal(Color{0, 0, 255, 1.0}))
			Expect(ParseCSSColor(`hsl(-240,100%,50%)`)).To(Equal(Color{0, 255, 0, 1.0}))
			Expect(ParseCSSColor(`hsl(840,100%,50%)`)).To(Equal(Color{0, 255, 0, 1.0}))
			Expect(ParseCSSColor(`hsl(840,100%,50%)`)).To(Equal(Color{0, 255, 0, 1.0}))
			Expect(ParseCSSColor(`hsl(-840,100%,50%)`)).To(Equal(Color{0, 0, 255, 1.0}))
		})

		It("should clamp the saturation and lightness correctly", func() {
			Expect(ParseCSSColor(`hsl(0,110%,50%)`)).To(Equal(Color{255, 0, 0, 1.0}))
			Expect(ParseCSSColor(`hsl(0,-120%,50%)`)).To(Equal(Color{127, 127, 127, 1.0}))
			Expect(ParseCSSColor(`hsl(0,100%,-100%)`)).To(Equal(Color{0, 0, 0, 1.0}))
			Expect(ParseCSSColor(`hsl(0,100%,200%)`)).To(Equal(Color{255, 255, 255, 1.0}))
		})
	})

	Context("when passed an hsla color", func() {
		It("should parse the color correctly", func() {
			Expect(ParseCSSColor(`hsla(0,0%,0%,1)`)).To(Equal(Color{0, 0, 0, 1.0}))
			Expect(ParseCSSColor(`hsla(0,100%,50%,0.5)`)).To(Equal(Color{255, 0, 0, 0.5}))
			Expect(ParseCSSColor(`hsla(120,100%,50%,0.2)`)).To(Equal(Color{0, 255, 0, 0.2}))
			Expect(ParseCSSColor(`hsla(240,100%,50%,0.7)`)).To(Equal(Color{0, 0, 255, 0.7}))
			Expect(ParseCSSColor(` hsla(0,100%,100%, 0) `)).To(Equal(Color{255, 255, 255, 0.0}))
		})

		It("should clamp the alpha value", func() {
			Expect(ParseCSSColor(`hsla(0,110%,50%,-0.1)`)).To(Equal(Color{255, 0, 0, 0.0}))
			Expect(ParseCSSColor(`hsla(0,-120%,50%,1.1)`)).To(Equal(Color{127, 127, 127, 1.0}))
		})
	})

	Context("when passed a color keyword", func() {
		It("should return the color for that keyword", func() {
			Expect(ParseCSSColor("aliceblue")).To(Equal(Color{240, 248, 255, 1.0}))
			Expect(ParseCSSColor("blue")).To(Equal(Color{0, 0, 255, 1.0}))
			Expect(ParseCSSColor("  coral ")).To(Equal(Color{255, 127, 80, 1.0}))
			Expect(ParseCSSColor("lightgoldenrodyellow")).To(Equal(Color{250, 250, 210, 1.0}))
		})
	})

	Context("cases that don't parse", func() {
		assertFailed := func(colorString string) {
			color, err := ParseCSSColor(colorString)
			ExpectWithOffset(1, color).To(BeZero())
			ExpectWithOffset(1, err).To(HaveOccurred())
		}

		Describe("some invalid layouts", func() {
			It("should error", func() {
				assertFailed("#0")
				assertFailed("#00")
				assertFailed("#-1-11111111111111111111111111111")
				assertFailed("#0000")
				assertFailed("#00000")
				assertFailed("#0000000")
				assertFailed("#GGG")
				assertFailed("#ggg")
				assertFailed("#aaaaag")
				assertFailed("rgb(1,2)")
				assertFailed("rgb(1.2,1.7,2.3)")
				assertFailed(`rgb(1,2,2.3%)`)
				assertFailed("rgb(1,2,3,4)")
				assertFailed(`rgb(1%,2%)`)
				assertFailed(`rgb(1%,2%,3%,4%)`)
				assertFailed("rgba(1,2,3)")
				assertFailed("rgba(1,2,3,5,6)")
				assertFailed(`rgba(1%,2%,3%,5%,6%)`)
				assertFailed(`hsl(120,5%,10)`)
				assertFailed(`hsl(120,5,10%)`)
				assertFailed(`hsl(120%,5%,10%)`)
				assertFailed(`hsl(120,5%,10%,1)`)
				assertFailed(`hsla(120,5%,10%)`)
				assertFailed(`hsla(120,5%,10%,1,1)`)
				assertFailed(`hsla(120,5%,10%,1%)`)
			})
		})

		Describe("numbers that aren't numbers", func() {
			It("should error", func() {
				invalidInteger := "11111111111111111111111111111111111"
				invalidFloat := "1.0.0"

				assertFailed("rgb(" + invalidInteger + ",2,3)")
				assertFailed("rgb(1," + invalidInteger + ",3)")
				assertFailed("rgb(1,2," + invalidInteger + ")")

				assertFailed("rgb(" + invalidFloat + "%,2%,3%)")
				assertFailed("rgb(1%," + invalidFloat + "%,3%)")
				assertFailed("rgb(1%,2%," + invalidFloat + "%)")

				assertFailed("rgba(" + invalidInteger + ",2,3,1.0)")
				assertFailed("rgba(1," + invalidInteger + ",3,1.0)")
				assertFailed("rgba(1,2," + invalidInteger + ",1.0)")
				assertFailed("rgba(1,2,3," + invalidFloat + ")")

				assertFailed("rgba(" + invalidFloat + "%,2%,3%,1.0)")
				assertFailed("rgba(1%," + invalidFloat + "%,3%,1.0)")
				assertFailed("rgba(1%,2%," + invalidFloat + "%,1.0)")
				assertFailed("rgba(1%,2%,3%," + invalidFloat + ")")

				assertFailed("hsl(" + invalidInteger + ",2%,3%)")
				assertFailed("hsl(1," + invalidFloat + "%,3%)")
				assertFailed("hsl(1,2%," + invalidFloat + "%)")

				assertFailed("hsla(" + invalidInteger + ",2%,3%,1.0)")
				assertFailed("hsla(1," + invalidFloat + "%,3%,1.0)")
				assertFailed("hsla(1,2%," + invalidFloat + "%,1.0)")
				assertFailed("hsla(1,2%,3%," + invalidFloat + ")")
			})
		})
	})
})
