#version 330

out vec4 fragColor;

uniform float time;
uniform vec2 resolution;

#define TAU 6.28318530718
#define PI 3.14159265359
#define MAX_ITER 12

void main() {
    float gamma     = 3.25;
    float speed     = 0.35;
    float scale     = 0.7;
    float brightness= 0.3;
    float contrast  = 1.25;

    float t1 = time * speed;

    // Safe normalize to [0,1] and apply scale
    vec2 res = max(resolution, vec2(1.0));
    vec2 uv  = (gl_FragCoord.xy / res) * scale;

    vec2 p = mod(uv * TAU, TAU) - 250.0;
    vec2 i = p;
    float c = 1.0;
    float inten = 0.005;

    for (int n = 0; n < MAX_ITER; n++) {
        float t = t1 * (1.0 - (3.5 / float(n + 1)));
        i = p + vec2(cos(t - i.x) + sin(t + i.y),
                     sin(t - i.y) + cos(t + i.x));
        c += 1.0 / length(vec2(
            p.x / (sin(i.x + t) / inten),
            p.y / (cos(i.y + t) / inten)
        ));
    }

    c /= float(MAX_ITER);
    c  = 1.17 - pow(c, 1.4);

    float value = pow(abs(c), 8.0);
    value += brightness;
    value  = mix(0.5, value, contrast);

    vec3 color = vec3(value);

    // Gamma correction
    fragColor = vec4(pow(color, vec3(gamma)), 1.0);
}
