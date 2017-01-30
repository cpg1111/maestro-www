'use strict';

const gulp = require('gulp');
const browserify = require('browserify');
const del = require('del');
const gulpUtil = require('gulp-util');
const source = require('vinyl-source-stream');
const buffer = require('vinyl-buffer');
const sourcemaps = require('gulp-sourcemaps');
const sass = require('gulp-ruby-sass');
const minifyCss = require('gulp-clean-css');
const htmlLint = require('gulp-html5-lint');

gulp.task('clean-scripts', ()=>del(['./build/js/**/*.js', './build/**/*.js.map']));

gulp.task('clean-styles', ()=>del(['./build/styles/**/*.css', './build/**/*.css.map']));

gulp.task('clean-templates', ()=>del(['./build/template/**/*.html']));

gulp.task('clean', ['clean-scripts', 'clean-styles', 'clean-template']);
    
var b = browserify({
    entries: ['./js/index.js'],
    extensions: ['.js'],
    debug: process.env['NODE_ENV'] != 'production'
});

b.transform('babelify', {
    presets: ['es2015', 'react', 'env']
});

function bundleScripts(){
    return b.bundle()
        .on('error', gulpUtil.log.bind(gulpUtil, 'Browserify Error'))
        .pipe(source('bundle.js'))
        .pipe(buffer())
        .pipe(sourcemaps.init({
            loadMaps: process.env['NODE_ENV'] != 'production'
        }))
        .pipe(sourcemaps.write('./maps'))
        .pipe(gulp.dest('./build/js'));
}

b.on('update', bundleScripts);

b.on('log', gulpUtil.log);

gulp.task('scripts', bundleScripts);

gulp.task('styles', ()=>{
    return gulp.src('./styles/**/index.sass')
        .pipe(sass().on('error', gulpUtil.log))
        .pipe(minifyCss())
        .pipe(gulp.dest('./build/styles'));
});

gulp.task('templates', ()=>{
    return gulp.src('../templates/index.html')
        .pipe(htmlLint())
        .pipe(gulp.dest('./build/template'));
});

gulp.task('default', ['scripts', 'styles', 'templates']);

gulp.task('watch-scripts', ()=>{
    watchify(b);
    bundleScripts();
});

gulp.task('watch-styles', ()=>gulp.watch('./styles/**/*', ['styles']));

gulp.task('watch-templates', ()=>gulp.watch('./templates/**/*', ['templates']));

gulp.task('watch', ['watch-scripts', 'watch-styles', 'watch-templates']);
