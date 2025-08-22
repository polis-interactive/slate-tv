#version 330

out vec4 fragColor;

uniform float time;
uniform vec2 resolution;

float noise(vec2 pos, float evolve) {
    // Loop the evolution (over a very long period of time).
    float e = fract(evolve * 0.01);

    // Coordinates
    float cx = pos.x * e;
    float cy = pos.y * e;

    // Generate a "random" black or white value
    return fract(23.0 * fract(2.0 / fract(fract(cx * 2.4 / cy * 23.0 + pow(abs(cy / 22.4), 3.3)) * fract(cx * evolve / pow(abs(cy), 0.050)))));
}

vec3 rgb2hsb(in vec3 c){
    vec4 K = vec4(0.0, -1.0/3.0, 2.0/3.0, -1.0);
    vec4 p = mix(vec4(c.bg, K.wz), vec4(c.gb, K.xy), step(c.b, c.g));
    vec4 q = mix(vec4(p.xyw, c.r),  vec4(c.r, p.yzx), step(p.x, c.r));
    float d = q.x - min(q.w, q.y);
    float e = 1.0e-10;
    return vec3(abs(q.z + (q.w - q.y) / (6.0*d + e)),
                d / (q.x + e),
                q.x);
}

vec3 hsb2rgb(in vec3 c){
    vec3 rgb = clamp(abs(mod(c.x*6.0 + vec3(0.0,4.0,2.0), 6.0) - 3.0) - 1.0,
                     0.0, 1.0);
    rgb = rgb*rgb*(3.0 - 2.0*rgb);
    return c.z * mix(vec3(1.0), rgb, c.y);
}

float rand(float n){ return fract(sin(n) * 43758.5453123); }

float noise(float p){
    float fl = floor(p);
    float fc = fract(p);
    return mix(rand(fl), rand(fl + 1.0), fc);
}

vec3 blacken(vec2 uv, float t) {
    return mix(vec3(1.0), vec3(0.0, 0.0, 1.0), noise(uv, t));
}

vec3 blacken(vec2 uv, float t, float intensity) {
    return mix(vec3(1.0), vec3(0.0, 0.0, 1.0), noise(uv, t)) * intensity;
}

void main() {
    float t = time / 2.0;

    vec2 res = max(resolution, vec2(1.0));
    vec2 uv  = gl_FragCoord.xy / res;

    vec3 color = vec3(0.0);

    float l = 1.0 / 7.0; // screen normalized to 1, divided by 7

    // 7 rectangles (top band)
    if (uv.y > 0.25) {
        if (uv.x >= l*0.0 && uv.x < l*1.0) color = vec3(1, 1, 1); // white
        if (uv.x >= l*1.0 && uv.x < l*2.0) color = vec3(1, 1, 0); // yellow
        if (uv.x >= l*2.0 && uv.x < l*3.0) color = vec3(0, 1, 1); // light blue
        if (uv.x >= l*3.0 && uv.x < l*4.0) color = vec3(0, 1, 0); // green
        if (uv.x >= l*4.0 && uv.x < l*5.0) color = vec3(1, 0, 1); // purple
        if (uv.x >= l*5.0 && uv.x < l*6.0) color = vec3(1, 0, 0); // red
        if (uv.x >= l*6.0 && uv.x < l*7.0) color = vec3(0, 0, 1); // blue
    }

    // 7 smaller rectangles (middle band)
    if (uv.y >= 0.2 && uv.y <= 0.25) {
        if (uv.x >= l*0.0 && uv.x < l*1.0) color = vec3(0, 0, 1);     // blue
        if (uv.x >= l*1.0 && uv.x < l*2.0) color = blacken(uv, t);    // noisy blue/black
        if (uv.x >= l*2.0 && uv.x < l*3.0) color = vec3(1, 0, 1);     // pink
        if (uv.x >= l*3.0 && uv.x < l*4.0) color = blacken(uv, t);    // noisy
        if (uv.x >= l*4.0 && uv.x < l*5.0) color = vec3(0, 1, 1);     // light blue
        if (uv.x >= l*5.0 && uv.x < l*6.0) color = blacken(uv, t);    // noisy
        if (uv.x >= l*6.0 && uv.x < l*7.0) color = vec3(1, 1, 1);     // white
    }

    // 6 squares (bottom band)
    float l2 = 1.0 / 6.0;
    if (uv.y < 0.2) {
        if (uv.x >= l2*0.0 && uv.x < l2*1.0) color = vec3(0, 0, 0.5); // dark blue
        if (uv.x >= l2*1.0 && uv.x < l2*2.0) color = vec3(1, 1, 1);   // white
        if (uv.x >= l2*2.0 && uv.x < l2*3.0) color = vec3(0.2, 0, 0.5); // dark purple
        // Gradient box
        if (uv.x >= l2*3.0 && uv.x <= l2*5.0) {
            float s = uv.x - l2*2.5;        // center gradient in the box
            color = blacken(uv, t, s);      // gradient via intensity
        }
        if (uv.x > l2*5.0 && uv.x < l2*6.0) color = blacken(uv, t); // noisy
    }

    fragColor = vec4(color, 1.0);
}
