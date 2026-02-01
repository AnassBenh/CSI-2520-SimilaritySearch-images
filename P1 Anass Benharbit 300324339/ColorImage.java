//Anass Benharbit 300324339


import java.awt.Color;
import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;
import javax.imageio.ImageIO;

public class ColorImage {
    private BufferedImage image;
    private int width;
    private int height;
    private int depth;

    public ColorImage(String filename) {
        try {
            File file = new File(filename);
            image = ImageIO.read(file);
            width = image.getWidth();
            height = image.getHeight();
            depth = 24; 
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public int[] getPixel(int x, int y) {
        int pixel = image.getRGB(x, y);
        int[] rgb = new int[3];
        rgb[0] = (pixel >> 16) & 0xFF; // Red
        rgb[1] = (pixel >> 8) & 0xFF;  // Green
        rgb[2] = pixel & 0xFF;         // Blue
        return rgb;
    }

    public void reduceColor(int d) {
        for (int y = 0; y < height; y++) {
            for (int x = 0; x < width; x++) {
                int[] rgb = getPixel(x, y);
                for (int i = 0; i < rgb.length; i++) {
                    rgb[i] = (rgb[i] / d) * d;
                }
                Color newColor = new Color(rgb[0], rgb[1], rgb[2]);
                image.setRGB(x, y, newColor.getRGB());
            }
        }
    }

    public int getWidth() {
        return width;
    }

    public int getHeight() {
        return height;
    }

    public int getDepth() {
        return depth;
    }
}