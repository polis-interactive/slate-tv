
#ifdef GL_ES

precision highp float;
#define IN varying
#define OUT out
#define TEXTURE texture2D

#else

#define IN in
#define OUT out
#define TEXTURE texture

#endif

uniform float time;
uniform vec2 resolution;
vec2 rot2D(vec2 v, float a)
{
    float cosa=cos(a);float sina=sin(a);
    mat2 rot = mat2(cosa,-sina,sina,cosa);
    return v*rot;
}

float sat(float a)
{
    return clamp(a,0.0,1.0);
}
float lenny(vec2 v)
{
    return abs(v.x)+abs(v.y);
}
int imod(int a, int b)
{
    int res = a/b;
    return a-(res*b);
}

// SDF operations
float _union(float a, float b)
{
    return min(a, b);
}
float _sub(float a, float b)
{
    return max(-a,b);
}

// SDF functions
float sdf_rect(vec2 uv, vec2 sz)
{
    vec2 r = abs(uv)-sz;
    return max(r.x,r.y);
}

float sdf_cir(vec2 uv, float r, float thick)
{
    float li = length(uv)-r;
    float lo = li-thick;
    return _sub(li,lo);
}

float sdf_plaind(vec2 uv)
{
    float left = uv.x;
    return _sub(left, sdf_cir(uv*vec2(0.6,1.0), -1.,1.075));
}

float sdf_v(vec2 uv)
{
    float res;
    float widthUp = 0.085;
    float pxu = 0.15;
    res = sdf_rect(uv-vec2(-pxu,0.089), vec2(widthUp,0.018));
    res = _union(res, sdf_rect(uv-vec2(pxu-0.02,0.089), vec2(widthUp-0.02,0.018)));
    res = _union(res, sdf_rect(rot2D(uv-vec2(-0.04,0.0005),1.1), vec2(0.1,0.025)));
    res = _union(res, sdf_rect(rot2D(uv-vec2(0.04,0.0005),-1.1), vec2(0.1,0.025)));
    return res;
}

float sdf_d(vec2 uv)
{
    float res = _sub(sdf_plaind(vec2(0.9,1.0)*uv+vec2(0.05,0.0)),sdf_plaind(uv*vec2(0.9,0.7)));

    res = _union(res, sdf_rect(uv-vec2(-0.015,0.089), vec2(0.03,0.018)));
    res = _union(res, sdf_rect(uv-vec2(-0.015,-0.089), vec2(0.03,0.018)));
    res = _union(res, sdf_rect(uv-vec2(-0.025,-0.02), vec2(0.02,0.07)));
    return res;
}

vec2 myPixel(vec2 uv, vec2 sz)
{
    vec2 tmp = uv / sz;

    uv.x = float(int(tmp.x));
    uv.y = float(int(tmp.y));
    return uv*sz;
}




vec3 rdrDvd(vec2 uv)
{
    float inC = float(length(uv*vec2(1.0,1.5))<0.2);
    float lVid = length(uv*vec2(0.9,4.5)*2.+vec2(0.0,0.49));
    float cVid = (1.0-sat((lVid-0.2)*5.0));
    return (vec3(1.0)*inC);
}
const float PI = 3.141592653927;
vec3 rdr(vec2 uv, float t)
{
    vec2 pos;
    vec3 col[4];

    col[0] = vec3(0.897,0.31,0.21);
    col[1] = vec3(0.897,0.31,0.21).xzy;
    col[2] = vec3(0.897,0.31,0.21).yzx;
    col[3] = vec3(0.1,0.897,0.21);

    pos.x = asin(sin(t))*0.5;
    bool isUpX = pos.x > asin(sin(t+0.01))*0.5;
    pos.y = asin(sin(2.0*t+PI))*0.55;
    bool isUpY = pos.y < asin(sin(2.0*t+PI+0.01))*0.25;
    int colIdx = int(isUpX)+(isUpY ? 1 : 2);
    return rdrDvd(uv-pos)*col[imod(colIdx,4)];
}


void main(void)
{

    float t = time / 2.0;
    vec2 uv = gl_FragCoord.xy / resolution.xy;
    vec3 color = vec3(0.0);

    vec2 center = vec2(0.5)*resolution.xy/resolution.xx;

    uv = uv - center;
    uv = uv *2.0;
    uv = myPixel(uv-vec2(2.0), vec2(0.01))+vec2(2.0);

    color = rdr(uv, t);

    gl_FragColor = vec4(color, 1.0);
}