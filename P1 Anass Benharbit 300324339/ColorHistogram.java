//Anass Benharbit 300324339


import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.PrintWriter;
import java.util.Arrays;
import java.util.Scanner;

public class ColorHistogram {
    private double[] histogram;

    public ColorHistogram(int d) {
        int numColors = (int) Math.pow(2, 3 * d);
        histogram = new double[numColors];
    }

    public ColorHistogram(String filename) {
        try {
            File file = new File(filename);
            Scanner scanner = new Scanner(file);
            histogram = new double[512]; 
            int i = 0;
            while (scanner.hasNextDouble()) {
                histogram[i++] = scanner.nextDouble();
            }
            scanner.close();
        } catch (FileNotFoundException e) {
            e.printStackTrace();
        }
    }

    public void setImage(ColorImage image) {
        Arrays.fill(histogram, 0);
        for (int y = 0; y < image.getHeight(); y++) {
            for (int x = 0; x < image.getWidth(); x++) {
                int[] rgb = image.getPixel(x, y);
                int index = rgb[0] / 32 + rgb[1] / 32 * 8 + rgb[2] / 32 * 64;
                histogram[index]++;
            }
        }
    }

    public double[] getHistogram() {
        return histogram;
    }

    public double compare(ColorHistogram hist) {
        double sum = 0;
        for (int i = 0; i < histogram.length; i++) {
            sum += Math.min(histogram[i], hist.histogram[i]);
        }
        return sum / Math.min(Arrays.stream(histogram).sum(), Arrays.stream(hist.histogram).sum());
    }

    public void save(String filename) {
        try {
            PrintWriter writer = new PrintWriter(filename);
            for (double value : histogram) {
                writer.println(value);
            }
            writer.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }


}