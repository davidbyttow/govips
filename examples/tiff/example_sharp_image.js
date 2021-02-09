// NOTE: Run from project root directory

const sharp = require("sharp");

let image = sharp("examples/tiff/input.jpg");

image.metadata().then((info) => {
	image
		// .rotate()
		.resize({
			width: info.width / 2,
			kernel: "lanczos3",
		})
		.tiff({
			quality: 100,
		})
		.toFile("examples/tiff/output-sharp.tiff");
});
